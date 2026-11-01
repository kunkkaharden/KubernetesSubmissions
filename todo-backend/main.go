package main

import (
	"cmp"
	"errors"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8080"

func main() {
	port := cmp.Or(os.Getenv("PORT"), defaultPort)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: newHandler(newTodoStore()),
	}

	log.Printf("Server started on port %s", port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server stopped: %v", err)
	}
}
