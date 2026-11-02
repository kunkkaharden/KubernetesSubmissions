package main

import (
	"fmt"
	"log"
	"net/http"
)

// ping answers with "pong <count>" and increments the counter.
func ping(c *counter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count, err := c.increment(r.Context())
		if err != nil {
			log.Printf("increment counter: %v", err)
			http.Error(w, "database unavailable", http.StatusInternalServerError)
			return
		}

		if _, err := fmt.Fprintf(w, "pong %d", count); err != nil {
			log.Printf("write ping response: %v", err)
		}
	}
}

// pings answers with the current count without incrementing it.
func pings(c *counter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count, err := c.current(r.Context())
		if err != nil {
			log.Printf("read counter: %v", err)
			http.Error(w, "database unavailable", http.StatusInternalServerError)
			return
		}

		if _, err := fmt.Fprintln(w, count); err != nil {
			log.Printf("write pings response: %v", err)
		}
	}
}
