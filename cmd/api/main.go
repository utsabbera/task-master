package main

import (
	"log"

	"github.com/utsabbera/task-master/api"
	_ "github.com/utsabbera/task-master/docs/swagger" // swaggo generated docs
)

// @title Task Master API
// @version 1.0
// @description API for managing tasks
// @host localhost:8080
// @BasePath /

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
