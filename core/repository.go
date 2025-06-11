package task

import (
	"errors"
	"sync"
	"time"

	"github.com/utsabbera/task-master/pkg/idgen"
)

var (
	// ErrTaskNotFound is returned when a task with the specified ID cannot be found
	ErrTaskNotFound = errors.New("task not found")
	// ErrInvalidTask is returned when an operation is performed on an invalid task
	ErrInvalidTask  = errors.New("invalid task")
)

// Repository defines the interface for task data storage operations
type Repository interface {
	// Create stores a new task in the repository and assigns it a unique ID
	Create(task *Task) error
	// FindByID retrieves a task by its ID
	FindByID(id string) (*Task, error)
	// FindAll returns all tasks stored in the repository
	FindAll() ([]*Task, error)
	// Update modifies an existing task in the repository
	Update(task *Task) error
	// Delete removes a task from the repository
	Delete(id string) error
}

// MemoryRepository is an in-memory implementation of Repository
// that stores tasks in a map and uses an ID generator for task IDs
type MemoryRepository struct {
	tasks     map[string]*Task
	generator idgen.Generator
	mu        sync.RWMutex
}

// NewMemoryRepository creates a new memory repository with the given ID generator
func NewMemoryRepository(idGenerator idgen.Generator) *MemoryRepository {
	return &MemoryRepository{
		tasks:     make(map[string]*Task),
		generator: idGenerator,
	}
}

// NewDefaultMemoryRepository creates a new MemoryRepository with a default sequential ID generator
func NewDefaultMemoryRepository() *MemoryRepository {
	generator := idgen.NewSequential("TASK-", 1, 6)
	return NewMemoryRepository(generator)
}

// Create stores a new task in memory and assigns it a unique ID and timestamps
func (r *MemoryRepository) Create(t *Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Set ID if not already set
	if t.ID == "" {
		t.ID = r.generator.Next()
	}
	
	// Set timestamps
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	r.tasks[t.ID] = t
	return nil
}

// FindByID retrieves a task by its ID from memory
// Returns ErrTaskNotFound if the task doesn't exist
func (r *MemoryRepository) FindByID(id string) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, exists := r.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}

	return t, nil
}

// FindAll returns all tasks stored in memory as a slice
func (r *MemoryRepository) FindAll() ([]*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// Update modifies an existing task in memory and updates the UpdatedAt timestamp
// Returns ErrTaskNotFound if the task doesn't exist
func (r *MemoryRepository) Update(t *Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[t.ID]; !exists {
		return ErrTaskNotFound
	}

	t.UpdatedAt = time.Now()
	r.tasks[t.ID] = t
	return nil
}

// Delete removes a task from memory by its ID
// Returns ErrTaskNotFound if the task doesn't exist
func (r *MemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(r.tasks, id)
	return nil
}
