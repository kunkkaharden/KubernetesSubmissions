package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type config struct {
	port        int
	databaseURL string
}

func loadConfig() (config, error) {
	port, err := positivePort("PORT")
	if err != nil {
		return config{}, err
	}
	dbPort, err := positivePort("DB_PORT")
	if err != nil {
		return config{}, err
	}

	dbHost, err := requiredEnv("DB_HOST")
	if err != nil {
		return config{}, err
	}
	dbUser, err := requiredEnv("POSTGRES_USER")
	if err != nil {
		return config{}, err
	}
	dbPassword, err := requiredEnv("POSTGRES_PASSWORD")
	if err != nil {
		return config{}, err
	}
	dbName, err := requiredEnv("POSTGRES_DB")
	if err != nil {
		return config{}, err
	}
	sslMode, err := requiredEnv("DB_SSLMODE")
	if err != nil {
		return config{}, err
	}

	connection := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbUser, dbPassword),
		Host:   fmt.Sprintf("%s:%d", dbHost, dbPort),
		Path:   "/" + dbName,
	}
	query := connection.Query()
	query.Set("sslmode", sslMode)
	connection.RawQuery = query.Encode()

	return config{
		port:        port,
		databaseURL: connection.String(),
	}, nil
}

func requiredEnv(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("%s is required", name)
	}
	return value, nil
}

func positivePort(name string) (int, error) {
	value, err := requiredEnv(name)
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(value)
	if err != nil || port < 1 || port > 65535 {
		return 0, fmt.Errorf("%s must be between 1 and 65535", name)
	}
	return port, nil
}
