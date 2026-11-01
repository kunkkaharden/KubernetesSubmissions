package main

import (
	"cmp"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort        = "8080"
	defaultPingPongURL = "http://ping-pong-svc:3000/pings"
)

func runServer(filePath string) {
	mux := http.NewServeMux()
	mux.Handle("GET /{$}", logOutputHandler{
		filePath:    filePath,
		pingPongURL: cmp.Or(os.Getenv("PINGPONG_URL"), defaultPingPongURL),
		client:      &http.Client{Timeout: 2 * time.Second},
	})

	port := cmp.Or(os.Getenv("PORT"), defaultPort)
	log.Printf("HTTP server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
