package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	defaultPort = "8080"
)

func pingPong(counter *int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong %d", *counter)
		*counter++
	}
}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		log.Fatalf("invalid PORT %q: must be a number between 1 and 65535", port)
	}

	counter := new(int)

	http.HandleFunc("GET /{$}", pingPong(counter))
	http.HandleFunc("GET /pingpong", pingPong(counter))

	fmt.Printf("Server started in port %d\n", portNumber)
	if err := http.ListenAndServe(":"+strconv.Itoa(portNumber), nil); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
