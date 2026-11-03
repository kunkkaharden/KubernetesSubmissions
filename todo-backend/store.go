package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// todo is a single task. The JSON field names are part of the public API.
type todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// todoStore keeps todos in Postgres so they survive application restarts.
type todoStore struct {
	db *sql.DB
}

// connectTodoStore waits for the database and creates the todos table on first run.
func connectTodoStore(databaseURL string) (*todoStore, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	for attempt := 1; attempt <= 30; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = db.PingContext(ctx)
		if err == nil {
			_, err = db.ExecContext(ctx, `
				CREATE TABLE IF NOT EXISTS todos (
					id BIGSERIAL PRIMARY KEY,
					text VARCHAR(140) NOT NULL,
					done BOOLEAN NOT NULL DEFAULT FALSE
				);
			`)
		}
		cancel()

		if err == nil {
			return &todoStore{db: db}, nil
		}

		log.Printf("database not ready (attempt %d/30): %v", attempt, err)
		time.Sleep(2 * time.Second)
	}

	db.Close()
	return nil, fmt.Errorf("connect to database: %w", err)
}

// list returns all todos ordered by ID.
func (s *todoStore) list(ctx context.Context) ([]todo, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, text, done FROM todos ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]todo, 0)
	for rows.Next() {
		var t todo
		if err := rows.Scan(&t.ID, &t.Text, &t.Done); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// add stores a new todo and returns it with its assigned ID.
func (s *todoStore) add(ctx context.Context, text string) (todo, error) {
	var t todo
	err := s.db.QueryRowContext(ctx,
		"INSERT INTO todos (text) VALUES ($1) RETURNING id, text, done",
		text,
	).Scan(&t.ID, &t.Text, &t.Done)
	return t, err
}
