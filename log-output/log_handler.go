package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// logOutputHandler returns the contents of the log file followed by the
// current Ping-pong count, which it asks for over HTTP.
type logOutputHandler struct {
	filePath    string
	pingPongURL string
	client      *http.Client
}

func (h logOutputHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile(h.filePath)
	if err != nil {
		log.Printf("read log file %q: %v", h.filePath, err)
		http.Error(w, "log file is not available", http.StatusInternalServerError)
		return
	}

	count, err := fetchPingPongCount(r.Context(), h.client, h.pingPongURL)
	if err != nil {
		log.Printf("fetch ping-pong count from %q: %v", h.pingPongURL, err)
		http.Error(w, "ping-pong counter is not available", http.StatusBadGateway)
		return
	}

	if len(content) > 0 && content[len(content)-1] != '\n' {
		content = append(content, '\n')
	}
	if _, err := fmt.Fprintf(w, "%sPing / Pongs: %d\n", content, count); err != nil {
		log.Printf("write response: %v", err)
	}
}
