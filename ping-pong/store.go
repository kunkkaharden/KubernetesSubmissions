package main

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"
)

const (
	defaultDBHost     = "localhost"
	defaultDBPort     = "5432"
	defaultDBUser     = "postgres"
	defaultDBPassword = "postgres"
	defaultDBName     = "pingpong"
)

// databaseURL builds the Postgres connection string. DATABASE_URL takes
// precedence; otherwise the connection is built from the standard
// POSTGRES_* and DB_* environment variables.
func databaseURL() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}

	connection := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cmp.Or(os.Getenv("POSTGRES_USER"), defaultDBUser), cmp.Or(os.Getenv("POSTGRES_PASSWORD"), defaultDBPassword)),
		Host:   cmp.Or(os.Getenv("DB_HOST"), defaultDBHost) + ":" + cmp.Or(os.Getenv("DB_PORT"), defaultDBPort),
		Path:   "/" + cmp.Or(os.Getenv("POSTGRES_DB"), defaultDBName),
	}
	query := connection.Query()
	query.Set("sslmode", cmp.Or(os.Getenv("DB_SSLMODE"), "disable"))
	connection.RawQuery = query.Encode()

	return connection.String()
}

// counter keeps the ping count in Postgres so it survives restarts.
type counter struct {
	db *sql.DB
}

// connectCounter waits for the database to accept connections and creates the
// counter table on first run.
func connectCounter(databaseURL string) (*counter, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	for attempt := 1; attempt <= 30; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = db.PingContext(ctx)
		if err == nil {
			_, err = db.ExecContext(ctx, `
				CREATE TABLE IF NOT EXISTS counter (
					id INTEGER PRIMARY KEY,
					value BIGINT NOT NULL
				);
				INSERT INTO counter (id, value) VALUES (1, 0)
				ON CONFLICT (id) DO NOTHING;
			`)
		}
		cancel()

		if err == nil {
			return &counter{db: db}, nil
		}

		log.Printf("database not ready (attempt %d/30): %v", attempt, err)
		time.Sleep(2 * time.Second)
	}

	db.Close()
	return nil, fmt.Errorf("connect to database: %w", err)
}

// increment adds one to the stored count and returns the new value.
func (c *counter) increment(ctx context.Context) (int64, error) {
	var count int64
	err := c.db.QueryRowContext(ctx,
		"UPDATE counter SET value = value + 1 WHERE id = 1 RETURNING value",
	).Scan(&count)
	return count, err
}

// current returns the stored count without changing it.
func (c *counter) current(ctx context.Context) (int64, error) {
	var count int64
	err := c.db.QueryRowContext(ctx, "SELECT value FROM counter WHERE id = 1").Scan(&count)
	return count, err
}
