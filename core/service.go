package task

import (
	"fmt"
	"time"
)

// Service provides business logic for task management operations
type Service struct {
	repo Repository
}

// NewService creates a new task service with the provided repository
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateTask creates a new task with the provided details
// Returns an error if the title is empty or if there's an issue with the repository
func (s *Service) CreateTask(title, description string, priority Priority, dueDate *time.Time) (*Task, error) {
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

// GetTask retrieves a task by its ID
// Returns an error if the ID is empty or if the task cannot be found
func (s *Service) GetTask(id string) (*Task, error) {
	if id == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("finding task: %w", err)
	}
	
	return task, nil
}

// GetAllTasks retrieves all tasks from the repository
func (s *Service) GetAllTasks() ([]*Task, error) {
	tasks, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("retrieving tasks: %w", err)
	}
	
	return tasks, nil
}

// UpdateTask updates an existing task in the repository
// Returns ErrInvalidTask if the task is nil
func (s *Service) UpdateTask(task *Task) error {
	if task == nil {
		return ErrInvalidTask
	}

	err := s.repo.Update(task)
	if err != nil {
		return fmt.Errorf("updating task: %w", err)
	}
	
	return nil
}

// DeleteTask removes a task from the repository by its ID
// Returns an error if the ID is empty or if the task cannot be found
func (s *Service) DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	
	err := s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("deleting task: %w", err)
	}
	
	return nil
}

// UpdateTaskStatus changes the status of a task with the specified ID
// Returns an error if the task cannot be found or updated
func (s *Service) UpdateTaskStatus(id string, status Status) error {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("finding task: %w", err)
	}
	
	task.Status = status
	
	err = s.repo.Update(task)
	if err != nil {
		return fmt.Errorf("updating task status: %w", err)
	}
	
	return nil
}
