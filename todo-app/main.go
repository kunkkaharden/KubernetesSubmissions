package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		log.Fatalf("invalid PORT %q: must be a number between 1 and 65535", port)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintln(w, "<!doctype html><html lang=\"en\"><head><meta charset=\"utf-8\"><title>Todo App</title></head><body><h1>Todo App</h1><p>Welcome!</p></body></html>")
	})

	fmt.Printf("Server started in port %d\n", portNumber)
	if err := http.ListenAndServe(":"+strconv.Itoa(portNumber), nil); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
