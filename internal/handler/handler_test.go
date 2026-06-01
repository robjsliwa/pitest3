package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/robjsliwa/pitest3/internal/storage"
)

func TestCreateTodoReturns201(t *testing.T) {
	store := storage.NewStore("")
	h := NewHandler(store)

	body := map[string]string{"title": "Buy milk"}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	http.HandlerFunc(h.CreateTodo).ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}

	var response map[string]map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	todo, ok := response["todo"]
	if !ok {
		t.Fatal("expected 'todo' key in response")
	}

	if todo["title"] != "Buy milk" {
		t.Errorf("expected title 'Buy milk', got %v", todo["title"])
	}

	if id, ok := todo["id"].(float64); !ok || int(id) != 1 {
		t.Errorf("expected id 1, got %v", todo["id"])
	}
}

func TestCreateTodoReturns400ForMissingTitle(t *testing.T) {
	store := storage.NewStore("")
	h := NewHandler(store)

	body := map[string]string{"title": ""}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	http.HandlerFunc(h.CreateTodo).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestCreateTodoReturns400ForWhitespaceTitle(t *testing.T) {
	store := storage.NewStore("")
	h := NewHandler(store)

	body := map[string]string{"title": "   "}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	http.HandlerFunc(h.CreateTodo).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestCreateTodoReturns400ForInvalidDueDate(t *testing.T) {
	store := storage.NewStore("")
	h := NewHandler(store)

	body := map[string]string{"title": "Valid", "due_date": "banana"}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	http.HandlerFunc(h.CreateTodo).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}
