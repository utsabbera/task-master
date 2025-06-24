package api

import (
	"net/http"

	assistant1 "github.com/utsabbera/task-master/core/assistant"
	"github.com/utsabbera/task-master/core/task"
	"github.com/utsabbera/task-master/pkg/assistant"
	"github.com/utsabbera/task-master/pkg/idgen"
	"github.com/utsabbera/task-master/pkg/middleware"
	"github.com/utsabbera/task-master/pkg/util"
)

// ServerConfig holds the configuration for the API server.
type ServerConfig struct {
	Addr      string
	Assistant assistant.Config
}

// NewServer returns a configured http.Server for the API.
func NewServer(cfg ServerConfig) *http.Server {
	addr := cfg.Addr
	if addr == "" {
		addr = ":8080"
	}

	repo := task.NewMemoryRepository()
	idGen := idgen.NewSequential("TASK-", 1, 6)
	clock := util.NewClock()
	taskService := task.NewService(repo, idGen, clock)
	assistant := assistant.NewClient(cfg.Assistant)
	assistantService := assistant1.NewService(taskService, assistant)
	handler := NewHandler(taskService, assistantService)

	middlewares := []middleware.Middleware{
		middleware.Log(),
	}

	return &http.Server{
		Addr:    addr,
		Handler: NewRouter(handler, middlewares...),
	}
}
