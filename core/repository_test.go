package task

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockGenerator struct {
	nextID int
	mu     sync.Mutex
}

func newMockGenerator(startID int) *mockGenerator {
	return &mockGenerator{nextID: startID}
}

func (g *mockGenerator) Next() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	id := g.nextID
	g.nextID++
	return string(rune('A' + id - 1))
}

func (g *mockGenerator) Current() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	return string(rune('A' + g.nextID - 1))
}

func (g *mockGenerator) Reset(id int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	g.nextID = id
}

func TestMemoryRepository_Create(t *testing.T) {
	t.Run("should assign ID to task when empty", func(t *testing.T) {
		generator := newMockGenerator(1)
		repo := NewMemoryRepository(generator)
		task := NewTask("Test Task", "Description", PriorityMedium, nil)

		assert.Empty(t, task.CreatedAt)
		assert.Empty(t, task.UpdatedAt)

		err := repo.Create(task)
		
		require.NoError(t, err)
		assert.Equal(t, "A", task.ID)

		assert.NotEmpty(t, task.CreatedAt)
		assert.NotEmpty(t, task.UpdatedAt)
		assert.Equal(t, task.CreatedAt, task.UpdatedAt)

		storedTask, err := repo.FindByID("A")
		require.NoError(t, err)
		assert.Equal(t, task, storedTask)
	})
	
	t.Run("should preserve existing ID", func(t *testing.T) {
		generator := newMockGenerator(1)
		repo := NewMemoryRepository(generator)
		task := NewTask("Test Task", "Description", PriorityMedium, nil)
		task.ID = "CUSTOM-ID"

		assert.Empty(t, task.CreatedAt)
		assert.Empty(t, task.UpdatedAt)

		err := repo.Create(task)
		
		require.NoError(t, err)
		assert.Equal(t, "CUSTOM-ID", task.ID)

		assert.NotEmpty(t, task.CreatedAt)
		assert.NotEmpty(t, task.UpdatedAt)
		
		storedTask, err := repo.FindByID("CUSTOM-ID")
		require.NoError(t, err)
		assert.Equal(t, task, storedTask)
	})
}

func TestMemoryRepository_FindByID(t *testing.T) {
	t.Run("should return task with matching ID", func(t *testing.T) {
		generator := newMockGenerator(1)
		repo := NewMemoryRepository(generator)
		task := NewTask("Test Task", "Description", PriorityMedium, nil)
		require.NoError(t, repo.Create(task))
		
		result, err := repo.FindByID(task.ID)
		
		require.NoError(t, err)
		assert.Equal(t, task, result)
	})
	
	t.Run("should return error when task not found", func(t *testing.T) {
		generator := newMockGenerator(1)
		repo := NewMemoryRepository(generator)
		
		result, err := repo.FindByID("non-existent")
		
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrTaskNotFound)
		assert.Nil(t, result)
	})
}

func TestMemoryRepository_FindAll(t *testing.T) {
	t.Run("should return all tasks", func(t *testing.T) {
		generator := newMockGenerator(1)
		repo := NewMemoryRepository(generator)
		task1 := NewTask("Task 1", "Description 1", PriorityLow, nil)
		task2 := NewTask("Task 2", "Description 2", PriorityHigh, nil)
		
		require.NoError(t, repo.Create(task1))
		require.NoError(t, repo.Create(task2))
		
		results, err := repo.FindAll()
		
		require.NoError(t, err)
		assert.Len(t, results, 2)
		
		foundTask1, foundTask2 := false, false
		for _, task := range results {
			if task.ID == task1.ID {
				foundTask1 = true
			}
			if task.ID == task2.ID {
				foundTask2 = true
			}
		}
		
		assert.True(t, foundTask1)
		assert.True(t, foundTask2)
	})

	t.Run("should return empty slice when no tasks", func(t *testing.T) {
		generator := newMockGenerator(1)
		repo := NewMemoryRepository(generator)
		
		results, err := repo.FindAll()
		
		require.NoError(t, err)
		assert.Empty(t, results)
	})
}

func TestMemoryRepository_Update(t *testing.T) {
	t.Run("should update existing task", func(t *testing.T) {
		generator := newMockGenerator(1)
		repo := NewMemoryRepository(generator)
		task := NewTask("Original Title", "Description", PriorityMedium, nil)
		require.NoError(t, repo.Create(task))

		originalCreatedAt := task.CreatedAt
		originalUpdatedAt := task.UpdatedAt
		
		time.Sleep(10 * time.Millisecond)
		
		task.Title = "Updated Title"
		task.Status = StatusInProgress
		
		err := repo.Update(task)
		
		require.NoError(t, err)
		
		updated, err := repo.FindByID(task.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", updated.Title)
		assert.Equal(t, StatusInProgress, updated.Status)

		assert.Equal(t, originalCreatedAt, updated.CreatedAt)
		assert.True(t, updated.UpdatedAt.After(originalUpdatedAt))
	})
	
	t.Run("should return error when task does not exist", func(t *testing.T) {	
		generator := newMockGenerator(1)
		repo := NewMemoryRepository(generator)
		task := NewTask("Original Title", "Description", PriorityMedium, nil)
		task.ID = "non-existent"
		
		err := repo.Update(task)
		
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})
}

func TestMemoryRepository_Delete(t *testing.T) {
	t.Run("should delete existing task", func(t *testing.T) {
		
		generator := newMockGenerator(1)
		repo := NewMemoryRepository(generator)
		task := NewTask("Test Task", "Description", PriorityMedium, nil)
		require.NoError(t, repo.Create(task))
		taskID := task.ID

		err := repo.Delete(taskID)
		require.NoError(t, err)
		
		_, err = repo.FindByID(taskID)
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})
	
	t.Run("should return error when task does not exist", func(t *testing.T) {
		generator := newMockGenerator(1)
		repo := NewMemoryRepository(generator)
		
		err := repo.Delete("non-existent")
		
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})
}

func TestNewDefaultMemoryRepository(t *testing.T) {
	t.Run("should create repository with sequential generator", func(t *testing.T) {
		
		repo := NewDefaultMemoryRepository()
		require.NotNil(t, repo)

		task := NewTask("Test Task", "Description", PriorityMedium, nil)
		require.NoError(t, repo.Create(task))
		assert.Equal(t, "TASK-000001", task.ID)
	})
}
