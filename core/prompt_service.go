package core

//go:generate mockgen -destination=mock_prompt_service.go -package=core . PromptService

// PromptService defines the interface for processing natural language prompts
// for task management operations
type PromptService interface {
	// ProcessPrompt handles a natural language prompt and performs the appropriate task operation
	ProcessPrompt(prompt string) (string, error)
}

type promptService struct {
	taskService TaskService
}

// NewPromptService creates a new prompt service with the provided task service
func NewPromptService(taskService TaskService) PromptService {
	return &promptService{
		taskService: taskService,
	}
}

// ProcessPrompt handles a natural language prompt and delegates to the appropriate task service method
func (s *promptService) ProcessPrompt(prompt string) (string, error) {
	// This is a placeholder implementation that will be expanded in future commits
	// In a full implementation, this would:
	// 1. Parse the prompt using NLP to determine intent (create, update, delete, list, etc.)
	// 2. Extract relevant information (task details, IDs, etc.)
	// 3. Call the appropriate TaskService method
	// 4. Return a human-readable response message

	// For now, let's implement a very simple version that creates a task with the prompt as title
	task, err := s.taskService.Create(
		prompt, // Use prompt as the title
		"",     // Empty description
		nil,    // No priority
		nil,    // No due date
	)

	if err != nil {
		return "", err
	}

	return "Task created: " + task.ID, nil
}
