package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	defaultPort     = "8080"
	defaultCacheDir = "./data"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		log.Fatalf("invalid PORT %q: must be a number between 1 and 65535", port)
	}

	cacheDir := os.Getenv("IMAGE_CACHE_DIR")
	if cacheDir == "" {
		cacheDir = defaultCacheDir
	}
	cache := newImageCache(cacheDir, fetchImage)

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(portNumber),
		Handler: newHandler(cache),
	}

	log.Printf("Server started on port %d", portNumber)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server stopped: %v", err)
	}
}
