package main

import "sync"

// todo is a single task. The JSON field names are part of the public API.
type todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// todoStore keeps todos in memory. A restart of the process clears them.
type todoStore struct {
	mu     sync.RWMutex
	todos  []todo
	nextID int
}

func newTodoStore() *todoStore {
	return &todoStore{
		todos:  make([]todo, 0),
		nextID: 1,
	}
}

// list returns a copy of the todos so callers cannot modify the store.
func (s *todoStore) list() []todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]todo, len(s.todos))
	copy(todos, s.todos)
	return todos
}

// add stores a new todo and returns it with its assigned ID.
func (s *todoStore) add(text string) todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := todo{ID: s.nextID, Text: text}
	s.nextID++
	s.todos = append(s.todos, t)
	return t
}
