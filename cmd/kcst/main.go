package main

import (
	"log"

	"github.com/keircn/kcst/internal/server"
)

func main() {
	srv := server.New(":8080")
	log.Println("Server starting on :8080")
	log.Fatal(srv.Run())
}
