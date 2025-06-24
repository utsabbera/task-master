package main

import (
	"log"

	"github.com/utsabbera/task-master/api"
	_ "github.com/utsabbera/task-master/docs/swagger" // swaggo generated docs
	"github.com/utsabbera/task-master/pkg/assistant"
)

// @title Task Master
// @version 1.0
// @description API for managing tasks
// @host localhost:8080
// @BasePath /

func main() {
	cfg := api.ServerConfig{
		Addr: ":8080",
		Assistant: assistant.Config{
			BaseURL:        "http://localhost:11434/v1",
			Model:          "llama3.2",
			AppName:        "Task Master",
			AppDescription: "AI powered application for managing tasks",
		},
	}

	server := api.NewServer(cfg)

	log.Println("Starting server on", cfg.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
