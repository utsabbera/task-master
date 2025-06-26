package task

import (
	"fmt"

	"github.com/utsabbera/task-master/pkg/idgen"
	"github.com/utsabbera/task-master/pkg/util"
)

//go:generate mockgen -destination=service_mock.go -package=task . Service

// Service defines the interface for task management operations
type Service interface {
	// Create adds a new task with the specified fields.
	// The caller must set Title, Description, Priority,DueDate, and Status.
	// The method mutates the provided *Task and returns an error if creation fails.
	Create(task *Task) error

	// Get retrieves a task by its ID
	// Returns an error if the task cannot be found
	Get(id string) (*Task, error)

	// List retrieves all tasks from the repository
	List() ([]*Task, error)

	// Update updates an existing task with the provided fields in the update parameter.
	// Only non-zero fields in the update parameter will overwrite the corresponding fields in the existing task.
	// Returns an error if the update parameter is nil, the task cannot be found, or the update operation fails.
	Update(id string, patch *Task) (*Task, error)

	// Delete removes a task from the repository by its ID
	// Returns an error if the task cannot be found
	Delete(id string) error
}

type service struct {
	repo        Repository
	clock       util.Clock
	idGenerator idgen.Generator
}

// NewService creates a new instance of Service with the provided repository, ID generator, and clock.
func NewService(repo Repository, idGenerator idgen.Generator, time util.Clock) Service {
	return &service{
		repo:        repo,
		clock:       time,
		idGenerator: idGenerator,
	}
}

func (s *service) Create(task *Task) error {
	if task.Status == "" {
		task.Status = StatusNotStarted
	}

	task.ID = s.idGenerator.Next()
	now := s.clock.Now()
	task.UpdatedAt = now
	task.CreatedAt = now

	err := s.repo.Create(task)
	if err != nil {
		return fmt.Errorf("error creating task: %w", err)
	}

	return nil
}

func (s *service) Get(id string) (*Task, error) {
	task, err := s.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("error finding task: %w", err)
	}

	return task, nil
}

func (s *service) List() ([]*Task, error) {
	tasks, err := s.repo.List()
	if err != nil {
		return nil, fmt.Errorf("error listing tasks: %w", err)
	}

	return tasks, nil
}

func (s *service) Update(id string, patch *Task) (*Task, error) {
	task, err := s.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("error finding task: %w", err)
	}

	s.update(task, patch)

	err = s.repo.Update(task)
	if err != nil {
		return nil, fmt.Errorf("error updating task: %w", err)
	}

	return task, nil
}

func (s *service) update(task *Task, patch *Task) {
	if patch.Title != "" {
		task.Title = patch.Title
	}
	if patch.Description != "" {
		task.Description = patch.Description
	}
	if patch.Priority != nil {
		task.Priority = patch.Priority
	}
	if patch.DueDate != nil {
		task.DueDate = patch.DueDate
	}
	if patch.Status != "" {
		task.Status = patch.Status
	}

	task.UpdatedAt = s.clock.Now()
}

func (s *service) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("error deleting task: %w", err)
	}

	return nil
}
