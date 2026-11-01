package main

import (
	"cmp"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

const defaultPort = "8080"

type counter struct {
	mu    sync.Mutex
	count int
}

func (c *counter) increment() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.count++
	return c.count
}

func (c *counter) current() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.count
}

// ping answers with "pong <count>" and increments the counter.
func ping(c *counter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintf(w, "pong %d", c.increment()); err != nil {
			log.Printf("write ping response: %v", err)
		}
	}
}

// pings answers with the current count without incrementing it.
func pings(c *counter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintln(w, c.current()); err != nil {
			log.Printf("write pings response: %v", err)
		}
	}
}

func main() {
	c := &counter{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", ping(c))
	mux.HandleFunc("GET /pingpong", ping(c))
	mux.HandleFunc("GET /pings", pings(c))

	port := cmp.Or(os.Getenv("PORT"), defaultPort)
	log.Printf("server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
