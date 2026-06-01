package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
)

func TestServerRespondsToRequests(t *testing.T) {
	srv, listener := setUpServer(t)
	defer srv.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://localhost:%d/todos", port)

	// Send invalid body — expect 400, proving the server is up and routing.
	resp, err := http.Post(url, "application/json", bytes.NewReader([]byte("not json")))
	if err != nil {
		t.Fatalf("server did not respond: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestServerCreateTodo(t *testing.T) {
	srv, listener := setUpServer(t)
	defer srv.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://localhost:%d/todos", port)

	body, _ := json.Marshal(map[string]string{"title": "Buy milk"})
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	var result map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	todoRaw, ok := result["todo"]
	if !ok {
		t.Fatal("expected 'todo' key in response")
	}

	var todo map[string]interface{}
	if err := json.Unmarshal(todoRaw, &todo); err != nil {
		t.Fatalf("failed to parse todo: %v", err)
	}

	if todo["title"] != "Buy milk" {
		t.Errorf("expected title 'Buy milk', got %v", todo["title"])
	}

	if id, ok := todo["id"].(float64); !ok || int(id) != 1 {
		t.Errorf("expected id 1, got %v", todo["id"])
	}
}

func TestServerCreateTodoRejectsEmptyTitle(t *testing.T) {
	srv, listener := setUpServer(t)
	defer srv.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://localhost:%d/todos", port)

	body, _ := json.Marshal(map[string]string{"title": ""})
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func setUpServer(t *testing.T) (*http.Server, net.Listener) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	wireMux(mux, "")

	srv := &http.Server{Handler: mux}

	go func() {
		srv.Serve(listener)
	}()

	t.Cleanup(func() { srv.Close() })
	return srv, listener
}
