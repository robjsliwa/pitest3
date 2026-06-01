package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/robjsliwa/pitest3/internal/handler"
	"github.com/robjsliwa/pitest3/internal/storage"
)

func main() {
	port := flag.Int("port", 8080, "HTTP port to listen on")
	dataFile := flag.String("data-file", "", "Path to the data file")
	flag.Parse()

	if *dataFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal("cannot determine home directory: ", err)
		}
		*dataFile = home + "/.pitest3/data.json"
	}

	// Ensure the parent directory for the data file exists.
	dataDir := filepath.Dir(*dataFile)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal("failed to create data directory: ", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}

	mux := http.NewServeMux()
	wireMux(mux, *dataFile)

	srv := &http.Server{Handler: mux}

	// Graceful shutdown on SIGINT / SIGTERM.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		log.Println("shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("forced shutdown: ", err)
		}
	}()

	log.Printf("server listening on :%d", *port)
	if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	log.Println("server stopped")
}

func wireMux(mux *http.ServeMux, dataFile string) {
	store := storage.NewStore(dataFile)
	h := handler.NewHandler(store)

	mux.HandleFunc("/todos", h.CreateTodo)
}
