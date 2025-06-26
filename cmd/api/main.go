package main

import (
	"log/slog"
	"os"
	"path/filepath"

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
	// TODO: Load configuration from environment variables or config file
	cfg := api.ServerConfig{
		Addr: ":8080",
		Assistant: assistant.Config{
			BaseURL:        "https://api.groq.com/openai/v1",
			APIKey:         "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",
			Model:          "meta-llama/llama-4-scout-17b-16e-instruct",
			AppName:        "Task Master",
			AppDescription: "AI powered application for managing tasks",
		},
	}

	server := api.NewServer(cfg)

	slog.Info("starting server", "addr", cfg.Addr)

	if err := server.ListenAndServe(); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source, _ := a.Value.Any().(*slog.Source)
				if source != nil {
					if filepath, err := filepath.Rel(os.ExpandEnv("$PWD"), source.File); err == nil {
						source.File = filepath
					} else {
						panic(err)
					}
				}
			}

			return a
		},
	})))
}
