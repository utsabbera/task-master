package task

import (
	"errors"
	"sync"
)

var (
	// ErrTaskNotFound is returned when a task with the specified ID cannot be found
	ErrTaskNotFound = errors.New("task not found")
	// ErrInvalidTask is returned when an operation is performed on an invalid task
	ErrInvalidTask = errors.New("invalid task")
)

//go:generate mockgen -destination=repository_mock.go -package=task . Repository

// Repository defines the interface for task data storage operations
type Repository interface {
	// Create stores a new task in the repository and assigns it a unique ID
	Create(task *Task) error

	// Get retrieves a task by its ID
	// Returns ErrTaskNotFound if the task doesn't exist
	Get(id string) (*Task, error)

	// List returns all tasks stored in the repository
	List() ([]*Task, error)

	// Update modifies an existing task in the repository
	// Returns ErrTaskNotFound if the task doesn't exist
	Update(task *Task) error

	// Delete removes a task from the repository
	// Returns ErrTaskNotFound if the task doesn't exist
	Delete(id string) error
}

// MemoryRepository is an in-memory implementation of Repository
// that stores tasks in a map and uses an ID generator for task IDs
type MemoryRepository struct {
	tasks map[string]*Task
	mu    sync.RWMutex
}

// NewMemoryRepository creates a new memory repository with the given ID generator
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks: make(map[string]*Task),
	}
}

func (r *MemoryRepository) Create(t *Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if t.ID == "" {
		return ErrInvalidTask
	}

	r.tasks[t.ID] = t
	return nil
}

func (r *MemoryRepository) Get(id string) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, exists := r.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}

	return t, nil
}

func (r *MemoryRepository) List() ([]*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *MemoryRepository) Update(t *Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[t.ID]; !exists {
		return ErrTaskNotFound
	}

	r.tasks[t.ID] = t
	return nil
}

func (r *MemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(r.tasks, id)
	return nil
}
