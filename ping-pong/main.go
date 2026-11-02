package main

import (
	"cmp"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

const defaultPort = "8080"

func main() {
	c, err := connectCounter(databaseURL())
	if err != nil {
		log.Fatal(err)
	}
	defer c.db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", ping(c))
	mux.HandleFunc("GET /pingpong", ping(c))
	mux.HandleFunc("GET /pings", pings(c))

	port := cmp.Or(os.Getenv("PORT"), defaultPort)
	log.Printf("server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
