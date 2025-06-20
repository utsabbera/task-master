package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mock_idgen "github.com/utsabbera/task-master/pkg/idgen"
)

func TestMemoryRepository_Create(t *testing.T) {
	t.Run("should assign ID to task when empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGen := mock_idgen.NewMockGenerator(ctrl)
		mockGen.EXPECT().Next().Return("A").Times(1)

		repo := NewMemoryRepository(mockGen)
		priority := PriorityMedium
		due := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		task := NewTask("Test Task", "Description", &priority, &due)

		assert.Empty(t, task.CreatedAt)
		assert.Empty(t, task.UpdatedAt)

		err := repo.Create(task)

		require.NoError(t, err)
		assert.Equal(t, "A", task.ID)

		assert.NotEmpty(t, task.CreatedAt)
		assert.NotEmpty(t, task.UpdatedAt)
		assert.Equal(t, task.CreatedAt, task.UpdatedAt)

		storedTask, err := repo.Get("A")
		require.NoError(t, err)
		assert.Equal(t, task, storedTask)
	})

	t.Run("should preserve existing ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGen := mock_idgen.NewMockGenerator(ctrl)
		// Next() should not be called when ID is already set

		repo := NewMemoryRepository(mockGen)
		task := NewTask("Test Task", "Description", nil, nil)
		task.ID = "CUSTOM-ID"

		assert.Empty(t, task.CreatedAt)
		assert.Empty(t, task.UpdatedAt)

		err := repo.Create(task)

		require.NoError(t, err)
		assert.Equal(t, "CUSTOM-ID", task.ID)

		assert.NotEmpty(t, task.CreatedAt)
		assert.NotEmpty(t, task.UpdatedAt)

		storedTask, err := repo.Get("CUSTOM-ID")
		require.NoError(t, err)
		assert.Equal(t, task, storedTask)
	})
	t.Run("should create task without priority", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGen := mock_idgen.NewMockGenerator(ctrl)
		mockGen.EXPECT().Next().Return("B").Times(1)

		repo := NewMemoryRepository(mockGen)
		due := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		task := NewTask("No Priority", "desc", nil, &due)

		err := repo.Create(task)
		require.NoError(t, err)
		assert.Equal(t, "B", task.ID)
		storedTask, err := repo.Get("B")

		require.NoError(t, err)
		assert.Nil(t, storedTask.Priority)
	})

	t.Run("should create task without dueDate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGen := mock_idgen.NewMockGenerator(ctrl)
		mockGen.EXPECT().Next().Return("C").Times(1)

		repo := NewMemoryRepository(mockGen)
		priority := PriorityMedium
		task := NewTask("No DueDate", "desc", &priority, nil)

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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGen := mock_idgen.NewMockGenerator(ctrl)
		mockGen.EXPECT().Next().Return("A").Times(1)

		repo := NewMemoryRepository(mockGen)
		task := NewTask("Test Task", "Description", nil, nil)
		require.NoError(t, repo.Create(task))

		result, err := repo.Get(task.ID)

		require.NoError(t, err)
		assert.Equal(t, task, result)
	})

	t.Run("should return error when task not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGen := mock_idgen.NewMockGenerator(ctrl)

		repo := NewMemoryRepository(mockGen)

		result, err := repo.Get("non-existent")

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrTaskNotFound)
		assert.Nil(t, result)
	})
}

func TestMemoryRepository_List(t *testing.T) {
	t.Run("should return all tasks", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGen := mock_idgen.NewMockGenerator(ctrl)
		mockGen.EXPECT().Next().Return("A").Times(1)
		mockGen.EXPECT().Next().Return("B").Times(1)

		repo := NewMemoryRepository(mockGen)
		task1 := NewTask("Task 1", "Description 1", nil, nil)
		task2 := NewTask("Task 2", "Description 2", nil, nil)

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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGen := mock_idgen.NewMockGenerator(ctrl)

		repo := NewMemoryRepository(mockGen)

		results, err := repo.List()

		require.NoError(t, err)
		assert.Empty(t, results)
	})
}

func TestMemoryRepository_Update(t *testing.T) {
	t.Run("should update existing task", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGen := mock_idgen.NewMockGenerator(ctrl)
		mockGen.EXPECT().Next().Return("A").Times(1)

		repo := NewMemoryRepository(mockGen)
		task := NewTask("Original Title", "Description", nil, nil)
		require.NoError(t, repo.Create(task))

		originalCreatedAt := task.CreatedAt
		originalUpdatedAt := task.UpdatedAt

		time.Sleep(10 * time.Millisecond)

		task.Title = "Updated Title"
		task.Status = StatusInProgress

		err := repo.Update(task)

		require.NoError(t, err)

		updated, err := repo.Get(task.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", updated.Title)
		assert.Equal(t, StatusInProgress, updated.Status)

		assert.Equal(t, originalCreatedAt, updated.CreatedAt)
		assert.True(t, updated.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("should return error when task does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGen := mock_idgen.NewMockGenerator(ctrl)

		repo := NewMemoryRepository(mockGen)
		task := NewTask("Original Title", "Description", nil, nil)
		task.ID = "non-existent"

		err := repo.Update(task)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})
}

func TestMemoryRepository_Delete(t *testing.T) {
	t.Run("should delete existing task", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGen := mock_idgen.NewMockGenerator(ctrl)
		mockGen.EXPECT().Next().Return("A").Times(1)

		repo := NewMemoryRepository(mockGen)
		task := NewTask("Test Task", "Description", nil, nil)
		require.NoError(t, repo.Create(task))
		taskID := task.ID

		err := repo.Delete(taskID)
		require.NoError(t, err)

		_, err = repo.Get(taskID)
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})

	t.Run("should return error when task does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGen := mock_idgen.NewMockGenerator(ctrl)

		repo := NewMemoryRepository(mockGen)

		err := repo.Delete("non-existent")

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})
}

func TestNewDefaultMemoryRepository(t *testing.T) {
	t.Run("should create repository with sequential generator", func(t *testing.T) {

		repo := NewDefaultMemoryRepository()
		require.NotNil(t, repo)

		task := NewTask("Test Task", "Description", nil, nil)
		require.NoError(t, repo.Create(task))
		assert.Equal(t, "TASK-000001", task.ID)
	})
}
