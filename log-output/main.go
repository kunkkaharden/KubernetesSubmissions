package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const (
	defaultPort = "8080"
	defaultFile = "/data/log.txt"
	writePeriod = 5 * time.Second
)

func currentStatus(randomValue string) string {
	return fmt.Sprintf("[%s] %s", time.Now().UTC().Format(time.RFC3339), randomValue)
}

func fileHandler(filePath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("read log file %q: %v", filePath, err)
			http.Error(w, "log file is not available", http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(content); err != nil {
			log.Printf("write HTTP response: %v", err)
		}
	})
}

func runWriter(filePath string) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalf("open log file %q: %v", filePath, err)
	}
	defer file.Close()

	randomValue := uuid.NewString()
	writeLine := func() {
		if _, err := fmt.Fprintln(file, currentStatus(randomValue)); err != nil {
			log.Fatalf("write log file %q: %v", filePath, err)
		}
	}

	writeLine()
	ticker := time.NewTicker(writePeriod)
	defer ticker.Stop()
	for range ticker.C {
		writeLine()
	}
}

func runServer(filePath string) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		log.Fatalf("invalid PORT %q: must be a number between 1 and 65535", port)
	}

	mux := http.NewServeMux()
	mux.Handle("/", fileHandler(filePath))

	address := ":" + strconv.Itoa(portNumber)
	log.Printf("HTTP server started on port %d", portNumber)
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func main() {
	filePath := os.Getenv("LOG_FILE")
	if filePath == "" {
		filePath = defaultFile
	}

	switch os.Getenv("ROLE") {
	case "writer":
		runWriter(filePath)
	case "server":
		runServer(filePath)
	default:
		log.Fatal("ROLE must be either \"writer\" or \"server\"")
	}
}
