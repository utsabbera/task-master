package assistant

import (
	"context"

	"github.com/utsabbera/task-master/core/task"
	"github.com/utsabbera/task-master/pkg/assistant"
)

//go:generate mockgen -destination=service_mock.go -package=assistant . Service

// Service defines the interface for processing natural language messages
// for task management operations
type Service interface {
	// Chat handles a natural language message and performs the appropriate task operation
	Chat(ctx context.Context, message string) (string, error)
}

type service struct {
	task      task.Service
	assistant assistant.Client
}

// NewService creates a new assistant service with the provided task service
func NewService(taskService task.Service, assistantClient assistant.Client) Service {
	service := &service{
		task:      taskService,
		assistant: assistantClient,
	}

	service.assistant.Init()

	return service
}

// Chat handles a natural language message and delegates to the appropriate task service method
func (s *service) Chat(ctx context.Context, message string) (string, error) {
	// TODO: Implement the logic to process the message using the assistant client
	//
	// return s.assistant.Chat(ctx, message)

	// Placeholder implementation for creating a task from the message
	task := &task.Task{
		Title:       message,
		Description: "",
		Priority:    nil,
		DueDate:     nil,
	}

	if err := s.task.Create(task); err != nil {
		return "", err
	}

	return "Task created: " + task.ID, nil
}
