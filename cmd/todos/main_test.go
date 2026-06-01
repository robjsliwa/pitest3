package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/robjsliwa/pitest3/internal/model"
)

func TestCreateTodoSucceeds(t *testing.T) {
	// Mock server that responds like our real API.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/todos" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var input struct {
			Title string `json:"title"`
		}
		json.NewDecoder(r.Body).Decode(&input)

		todo := model.Todo{
			ID:    1,
			Title: input.Title,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]model.Todo{"todo": todo})
	}))
	defer srv.Close()

	result, err := createTodo(srv.URL, "Buy milk", "", "")
	if err != nil {
		t.Fatalf("createTodo failed: %v", err)
	}

	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}

	if result.Title != "Buy milk" {
		t.Errorf("expected title 'Buy milk', got %q", result.Title)
	}
}

func TestCreateTodoFailsWhenServerRejects(t *testing.T) {
	// Mock server that returns 400.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	_, err := createTodo(srv.URL, "", "", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFormatTodo(t *testing.T) {
	todo := model.Todo{
		ID:    1,
		Title: "Buy milk",
	}
	output := formatTodo(todo)

	// Should contain the key fields.
	if !strings.Contains(output, "ID:") || !strings.Contains(output, "Title:") {
		t.Errorf("formatTodo output missing expected fields: %q", output)
	}

	if !strings.Contains(output, "Buy milk") {
		t.Errorf("formatTodo output missing title: %q", output)
	}
}
