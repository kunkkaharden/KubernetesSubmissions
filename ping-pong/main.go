package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	defaultPort       = "8080"
	defaultCounterFile = "/data/pingpong.txt"
)

type counterStore struct {
	filePath string
	count    int
	mu       sync.Mutex
}

func newCounterStore(filePath string) (*counterStore, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("read counter file %q: %w", filePath, err)
		}
		if err := os.WriteFile(filePath, []byte("0\n"), 0o644); err != nil {
			return nil, fmt.Errorf("initialize counter file %q: %w", filePath, err)
		}
		content = []byte("0")
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil || count < 0 {
		return nil, fmt.Errorf("invalid counter in %q: %q", filePath, strings.TrimSpace(string(content)))
	}
	return &counterStore{filePath: filePath, count: count}, nil
}

func (s *counterStore) increment() (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := s.count + 1
	tempPath := s.filePath + ".tmp"
	if err := os.WriteFile(tempPath, []byte(strconv.Itoa(count)+"\n"), 0o644); err != nil {
		return 0, fmt.Errorf("write counter file %q: %w", tempPath, err)
	}
	if err := os.Rename(tempPath, s.filePath); err != nil {
		return 0, fmt.Errorf("replace counter file %q: %w", s.filePath, err)
	}
	s.count = count
	return count, nil
}

func pingPong(counter *counterStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count, err := counter.increment()
		if err != nil {
			log.Printf("increment ping-pong counter: %v", err)
			http.Error(w, "ping-pong counter is not available", http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprintf(w, "pong %d", count); err != nil {
			log.Printf("write HTTP response: %v", err)
		}
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

	counterFile := os.Getenv("COUNTER_FILE")
	if counterFile == "" {
		counterFile = defaultCounterFile
	}
	counter, err := newCounterStore(counterFile)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("GET /{$}", pingPong(counter))
	http.HandleFunc("GET /pingpong", pingPong(counter))

	fmt.Printf("Server started in port %d\n", portNumber)
	if err := http.ListenAndServe(":"+strconv.Itoa(portNumber), nil); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
