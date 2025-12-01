package main

import (
	"flag"
	"log"

	"github.com/keircn/kcst/internal/config"
	"github.com/keircn/kcst/internal/server"
)

func main() {
	configPath := flag.String("config", "config.toml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	log.Printf("Server starting on %s", cfg.Server.Address)
	log.Fatal(srv.Run())
}
