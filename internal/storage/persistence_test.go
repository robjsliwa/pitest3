package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robjsliwa/pitest3/internal/model"
)

func TestStorePersistsToJSONFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "todos.json")

	store := NewStore(file)

	todo, err := store.Create(model.Todo{
		Title:       "Buy milk",
		Description: "Organic, 2%",
		DueDate:     "2025-01-20",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if todo.ID != 1 {
		t.Errorf("expected ID 1, got %d", todo.ID)
	}

	// Verify file exists and contains the Todo.
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("data file not created: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("data file is empty")
	}

	// Verify the file contains the title.
	if !strings.Contains(string(data), "Buy milk") {
		t.Errorf("data file missing title: %s", string(data))
	}
}

func TestStoreLoadsFromJSONFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "todos.json")

	// Create a store and add a Todo.
	store1 := NewStore(file)
	_, err := store1.Create(model.Todo{Title: "Existing todo"})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Create a new store from the same file — should load existing data.
	store2 := NewStore(file)

	list := store2.List()
	if len(list) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(list))
	}

	if list[0].Title != "Existing todo" {
		t.Errorf("expected title 'Existing todo', got %q", list[0].Title)
	}

	// Next ID should continue from where we left off.
	next, err := store2.Create(model.Todo{Title: "New todo"})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if next.ID != 2 {
		t.Errorf("expected ID 2, got %d", next.ID)
	}
}

func TestStoreLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "nonexistent.json")

	// Should not fail — just start fresh.
	store := NewStore(file)

	list := store.List()
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d items", len(list))
	}
}
