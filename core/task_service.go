package core

import (
	"fmt"
	"time"
)

//go:generate mockgen -destination=mock_task_service.go -package=core . TaskService

// TaskService defines the interface for task management operations
type TaskService interface {
	// Create adds a new task with the specified title, description, priority, and optional due date
	Create(title, description string, priority *Priority, dueDate *time.Time) (*Task, error)

	// Get retrieves a task by its ID
	Get(id string) (*Task, error)

	// List returns all tasks stored in the repository
	List() ([]*Task, error)

	// Update modifies an existing task in the repository
	Update(task *Task) error

	// Delete removes a task from the repository by its ID
	Delete(id string) error
}

type taskService struct {
	repo Repository
}

// NewTaskService creates a new task service with the provided repository
func NewTaskService(repo Repository) TaskService {
	return &taskService{
		repo: repo,
	}
}

// Create creates a new task with the provided details
// Returns an error if the title is empty or if there's an issue with the repository
func (s *taskService) Create(title, description string, priority *Priority, dueDate *time.Time) (*Task, error) {
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}

	task := NewTask(title, description, priority, dueDate)

	err := s.repo.Create(task)
	if err != nil {
		return nil, fmt.Errorf("creating task: %w", err)
	}

	return task, nil
}

// Get retrieves a task by its ID
// Returns an error if the ID is empty or if the task cannot be found
func (s *taskService) Get(id string) (*Task, error) {
	if id == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}

	task, err := s.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("finding task: %w", err)
	}

	return task, nil
}

// List retrieves all tasks from the repository
func (s *taskService) List() ([]*Task, error) {
	tasks, err := s.repo.List()
	if err != nil {
		return nil, fmt.Errorf("retrieving tasks: %w", err)
	}

	return tasks, nil
}

// Update updates an existing task in the repository
// Returns ErrInvalidTask if the task is nil
func (s *taskService) Update(task *Task) error {
	if task == nil {
		return ErrInvalidTask
	}

	err := s.repo.Update(task)
	if err != nil {
		return fmt.Errorf("updating task: %w", err)
	}

	return nil
}

// Delete removes a task from the repository by its ID
// Returns an error if the ID is empty or if the task cannot be found
func (s *taskService) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("task ID cannot be empty")
	}

	err := s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("deleting task: %w", err)
	}

	return nil
}
