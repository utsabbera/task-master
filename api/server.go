package api

import (
	"net/http"

	"github.com/utsabbera/task-master/core"
	"github.com/utsabbera/task-master/pkg/middleware"
)

// ServerConfig holds the configuration for the API server.
type ServerConfig struct {
	Addr string
}

// NewServer returns a configured http.Server for the API.
func NewServer(cfg ServerConfig) *http.Server {
	addr := cfg.Addr
	if addr == "" {
		addr = ":8080"
	}

	repo := core.NewDefaultMemoryRepository()
	taskService := core.NewTaskService(repo)
	promptService := core.NewPromptService(taskService)
	handler := NewHandler(taskService, promptService)

	middlewares := []middleware.Middleware{
		middleware.Log(),
	}

	return &http.Server{
		Addr:    addr,
		Handler: NewRouter(handler, middlewares...),
	}
}
