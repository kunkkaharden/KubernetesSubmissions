package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"
)

const maxTodoLength = 140

type createTodoRequest struct {
	Text string `json:"text"`
}

func newHandler(store *todoStore) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", listTodos(store))
	mux.HandleFunc("POST /todos", createTodo(store))
	return mux
}

func listTodos(store *todoStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, store.list())
	}
}

func createTodo(store *todoStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		var request createTodoRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		text := strings.TrimSpace(request.Text)
		if text == "" {
			http.Error(w, "Todo text is required", http.StatusBadRequest)
			return
		}
		if utf8.RuneCountInString(text) > maxTodoLength {
			http.Error(w, "Todo text must be at most 140 characters", http.StatusBadRequest)
			return
		}

		writeJSON(w, http.StatusCreated, store.add(text))
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}
