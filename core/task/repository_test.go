package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryRepository_Create(t *testing.T) {
	t.Run("should return error when id is empty", func(t *testing.T) {
		repo := NewMemoryRepository()
		priority := PriorityMedium
		due := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		task := NewTask("Test Task", "Description", &priority, &due)

		err := repo.Create(task)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidTask)
	})

	t.Run("should preserve existing ID", func(t *testing.T) {
		repo := NewMemoryRepository()
		task := &Task{ID: "CUSTOM-ID", Title: "Test Task", Description: "Description"}

		err := repo.Create(task)
		assert.NoError(t, err)

		storedTask, err := repo.Get("CUSTOM-ID")
		require.NoError(t, err)
		assert.Equal(t, task, storedTask)
	})

	t.Run("should create task without priority", func(t *testing.T) {
		repo := NewMemoryRepository()
		due := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		task := &Task{ID: "B", Title: "No Priority", Description: "desc", DueDate: &due}

		err := repo.Create(task)
		require.NoError(t, err)
		storedTask, err := repo.Get("B")

		require.NoError(t, err)
		assert.Nil(t, storedTask.Priority)
	})

	t.Run("should create task without dueDate", func(t *testing.T) {
		repo := NewMemoryRepository()
		priority := PriorityMedium
		task := &Task{ID: "C", Title: "No DueDate", Description: "desc", Priority: &priority}

		err := repo.Create(task)
		require.NoError(t, err)
		assert.Equal(t, "C", task.ID)

		storedTask, err := repo.Get("C")
		require.NoError(t, err)
		assert.Nil(t, storedTask.DueDate)
	})
}

func TestMemoryRepository_Get(t *testing.T) {
	t.Run("should return task with matching ID", func(t *testing.T) {
		repo := NewMemoryRepository()
		task := &Task{ID: "A", Title: "Test Task", Description: "Description"}

		require.NoError(t, repo.Create(task))

		result, err := repo.Get(task.ID)

		require.NoError(t, err)
		assert.Equal(t, task, result)
	})

	t.Run("should return error when task not found", func(t *testing.T) {
		repo := NewMemoryRepository()

		result, err := repo.Get("non-existent")

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrTaskNotFound)
		assert.Nil(t, result)
	})
}

func TestMemoryRepository_List(t *testing.T) {
	t.Run("should return all tasks", func(t *testing.T) {
		repo := NewMemoryRepository()
		task1 := &Task{ID: "task1-id", Title: "Task 1", Description: "Description 1"}
		task2 := &Task{ID: "task2-id", Title: "Task 2", Description: "Description 2"}

		require.NoError(t, repo.Create(task1))
		require.NoError(t, repo.Create(task2))

		results, err := repo.List()

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
		repo := NewMemoryRepository()

		results, err := repo.List()

		require.NoError(t, err)
		assert.Empty(t, results)
	})
}

func TestMemoryRepository_Update(t *testing.T) {
	t.Run("should update existing task", func(t *testing.T) {
		repo := NewMemoryRepository()
		task := &Task{ID: "task-id", Title: "Original Title", Description: "Description"}

		require.NoError(t, repo.Create(task))

		task.Title = "Updated Title"
		task.Status = StatusInProgress

		err := repo.Update(task)
		require.NoError(t, err)

		updated, err := repo.Get(task.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", updated.Title)
		assert.Equal(t, StatusInProgress, updated.Status)
	})

	t.Run("should return error when task does not exist", func(t *testing.T) {
		repo := NewMemoryRepository()
		task := &Task{ID: "non-existent", Title: "Original Title", Description: "Description"}

		err := repo.Update(task)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})
}

func TestMemoryRepository_Delete(t *testing.T) {
	t.Run("should delete existing task", func(t *testing.T) {
		repo := NewMemoryRepository()
		task := &Task{ID: "A", Title: "Test Task", Description: "Description"}

		require.NoError(t, repo.Create(task))

		err := repo.Delete(task.ID)
		require.NoError(t, err)

		_, err = repo.Get(task.ID)
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})

	t.Run("should return error when task does not exist", func(t *testing.T) {
		repo := NewMemoryRepository()

		err := repo.Delete("non-existent")

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})
}
