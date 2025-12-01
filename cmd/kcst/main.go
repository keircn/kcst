package main

import (
	"log"

	"github.com/keircn/kcst/internal/server"
)

func main() {
	srv, err := server.New(":8080", "./uploads")
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	log.Println("Server starting on :8080")
	log.Fatal(srv.Run())
}
