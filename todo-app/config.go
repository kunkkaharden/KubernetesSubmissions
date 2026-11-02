package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type config struct {
	port              int
	imageCacheDir     string
	imageCacheTTL     time.Duration
	imageURL          string
	imageFetchTimeout time.Duration
	maxImageSize      int64
	todoBackendURL    string
}

func loadConfig() (config, error) {
	port, err := positiveInt("PORT")
	if err != nil {
		return config{}, err
	}
	if port > 65535 {
		return config{}, fmt.Errorf("PORT must be between 1 and 65535")
	}

	cacheDir, err := requiredEnv("IMAGE_CACHE_DIR")
	if err != nil {
		return config{}, err
	}
	cacheTTL, err := positiveDuration("IMAGE_CACHE_TTL")
	if err != nil {
		return config{}, err
	}
	imageURL, err := requiredEnv("IMAGE_URL")
	if err != nil {
		return config{}, err
	}
	fetchTimeout, err := positiveDuration("IMAGE_FETCH_TIMEOUT")
	if err != nil {
		return config{}, err
	}
	maxImageSize, err := positiveInt64("MAX_IMAGE_SIZE_BYTES")
	if err != nil {
		return config{}, err
	}
	backendURL, err := requiredEnv("TODO_BACKEND_URL")
	if err != nil {
		return config{}, err
	}

	return config{
		port:              port,
		imageCacheDir:     cacheDir,
		imageCacheTTL:     cacheTTL,
		imageURL:          imageURL,
		imageFetchTimeout: fetchTimeout,
		maxImageSize:      maxImageSize,
		todoBackendURL:    backendURL,
	}, nil
}

func requiredEnv(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("%s is required", name)
	}
	return value, nil
}

func positiveInt(name string) (int, error) {
	value, err := requiredEnv(name)
	if err != nil {
		return 0, err
	}
	number, err := strconv.Atoi(value)
	if err != nil || number < 1 {
		return 0, fmt.Errorf("%s must be a positive integer", name)
	}
	return number, nil
}

func positiveInt64(name string) (int64, error) {
	value, err := requiredEnv(name)
	if err != nil {
		return 0, err
	}
	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil || number < 1 {
		return 0, fmt.Errorf("%s must be a positive integer", name)
	}
	return number, nil
}

func positiveDuration(name string) (time.Duration, error) {
	value, err := requiredEnv(name)
	if err != nil {
		return 0, err
	}
	duration, err := time.ParseDuration(value)
	if err != nil || duration <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration (for example, 30s or 10m)", name)
	}
	return duration, nil
}
