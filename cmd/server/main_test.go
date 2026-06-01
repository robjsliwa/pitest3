package main

import (
	"fmt"
	"net"
	"net/http"
	"testing"
)

func TestServerRespondsToRequests(t *testing.T) {
	srv, listener := setUpServer(t)
	defer srv.Close()

	port := listener.Addr().(*net.TCPAddr).Port

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/", port))
	if err != nil {
		t.Fatalf("server did not respond: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func setUpServer(t *testing.T) (*http.Server, net.Listener) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	go func() {
		srv.Serve(listener)
	}()

	t.Cleanup(func() { srv.Close() })
	return srv, listener
}
