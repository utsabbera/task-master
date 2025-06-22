package chat

import "github.com/utsabbera/task-master/core/task"

//go:generate mockgen -destination=service_mock.go -package=chat . Service

// Service defines the interface for processing natural language prompts
// for task management operations
type Service interface {
	// ProcessPrompt handles a natural language prompt and performs the appropriate task operation
	ProcessPrompt(prompt string) (string, error)
}

type service struct {
	taskService task.Service
}

// NewService creates a new prompt service with the provided task service
func NewService(taskService task.Service) Service {
	return &service{
		taskService: taskService,
	}
}

// ProcessPrompt handles a natural language prompt and delegates to the appropriate task service method
func (s *service) ProcessPrompt(prompt string) (string, error) {
	// This is a placeholder implementation that will be expanded in future commits
	// In a full implementation, this would:
	// 1. Parse the prompt using NLP to determine intent (create, update, delete, list, etc.)
	// 2. Extract relevant information (task details, IDs, etc.)
	// 3. Call the appropriate Service method
	// 4. Return a human-readable response message

	// For now, let's implement a very simple version that creates a task with the prompt as title

	task := &task.Task{
		Title:       prompt,
		Description: "",
		Priority:    nil,
		DueDate:     nil,
	}

	if err := s.taskService.Create(task); err != nil {
		return "", err
	}

	return "Task created: " + task.ID, nil
}
