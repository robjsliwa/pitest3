package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}

	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	log.Printf("server listening on :%d", *port)
	if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
