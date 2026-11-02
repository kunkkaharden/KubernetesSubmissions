package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
)

func main() {
	settings, err := loadConfig()
	if err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	cache := newImageCache(settings.imageCacheDir, settings.imageCacheTTL, func(ctx context.Context) ([]byte, error) {
		return fetchImage(ctx, settings.imageURL, settings.imageFetchTimeout, settings.maxImageSize)
	})

	todoProxy, err := newTodoProxy(settings.todoBackendURL)
	if err != nil {
		log.Fatalf("invalid TODO_BACKEND_URL %q: %v", settings.todoBackendURL, err)
	}

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(settings.port),
		Handler: newHandler(cache, todoProxy),
	}

	log.Printf("Server started on port %d", settings.port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server stopped: %v", err)
	}
}
