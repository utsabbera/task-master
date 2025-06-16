package main

import (
	"log"

	"github.com/utsabbera/task-master/api"
)

func main() {
	cfg := api.ServerConfig{
		Addr: ":8080",
	}

	server := api.NewServer(cfg)

	log.Println("Starting server on", cfg.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
