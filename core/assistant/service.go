package assistant

import (
	"context"

	"github.com/utsabbera/task-master/core/task"
	"github.com/utsabbera/task-master/pkg/assistant"
	"github.com/utsabbera/task-master/pkg/util"
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
func NewService(
	taskService task.Service,
	assistantClient assistant.Client,
	clock util.Clock,
) Service {
	service := &service{
		task:      taskService,
		assistant: assistantClient,
	}

	service.assistant.RegisterFunctions(
		NewDateTimeFunction(clock),
		NewCreateFunction(taskService),
		NewGetFunction(taskService),
		NewListFunction(taskService),
		NewUpdateFunction(taskService),
		NewDeleteFunction(taskService),
	)

	service.assistant.Init()

	return service
}

// Chat handles a natural language message and delegates to the appropriate task service method
func (s *service) Chat(ctx context.Context, message string) (string, error) {
	return s.assistant.Chat(ctx, message)
}
