package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

func main() {
	settings, err := loadConfig()
	if err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	store, err := connectTodoStore(settings.databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer store.db.Close()

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(settings.port),
		Handler: newHandler(store),
	}

	log.Printf("Server started on port %d", settings.port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server stopped: %v", err)
	}
}
