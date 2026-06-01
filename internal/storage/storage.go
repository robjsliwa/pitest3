package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/robjsliwa/pitest3/internal/model"
)

// Storage persists Todos in memory with JSON file backup.
type Storage struct {
	mu     sync.RWMutex
	todos  []model.Todo
	nextID int
	file   string
}

// NewStore creates a new Storage, loading from file if it exists.
func NewStore(file string) *Storage {
	s := &Storage{
		todos:  []model.Todo{},
		nextID: 1,
		file:   file,
	}
	s.load()
	return s
}

// load reads existing Todos from the JSON data file.
// Safe to call without holding s.mu because NewStore is always called once
// at startup before any goroutine accesses the store. If load() is ever
// called after the store is published, acquire s.mu.Lock() first.
func (s *Storage) load() {
	data, err := os.ReadFile(s.file)
	if err != nil {
		// File doesn't exist or can't be read — start fresh.
		return
	}

	var loaded []model.Todo
	if err := json.Unmarshal(data, &loaded); err != nil {
		return
	}

	s.todos = loaded
	// Calculate nextID from loaded data.
	for _, t := range loaded {
		if t.ID >= s.nextID {
			s.nextID = t.ID + 1
		}
	}
}

// save writes all Todos to the JSON data file atomically.
func (s *Storage) save() error {
	if s.file == "" {
		return nil
	}

	data, err := json.MarshalIndent(s.todos, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal todos: %w", err)
	}

	tmp := s.file + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}

	if err := os.Rename(tmp, s.file); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}

	return nil
}

// Create validates, assigns an ID, stores the Todo, and persists to disk.
func (s *Storage) Create(todo model.Todo) (model.Todo, error) {
	if err := todo.Validate(); err != nil {
		return model.Todo{}, fmt.Errorf("validate todo: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	todo.ID = s.nextID
	s.nextID++
	todo.CreatedAt = now
	todo.UpdatedAt = now

	s.todos = append(s.todos, todo)

	if err := s.save(); err != nil {
		return model.Todo{}, fmt.Errorf("save: %w", err)
	}

	return todo, nil
}

// GetByID returns a Todo by its ID.
func (s *Storage) GetByID(id int) (model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, t := range s.todos {
		if t.ID == id {
			return t, nil
		}
	}

	return model.Todo{}, errors.New("todo not found")
}

// List returns all Todos.
func (s *Storage) List() []model.Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.Todo, len(s.todos))
	copy(result, s.todos)
	return result
}
