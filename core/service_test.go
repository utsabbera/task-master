package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utsabbera/task-master/pkg/match"
	"go.uber.org/mock/gomock"
)

func TestService_Create(t *testing.T) {
	t.Run("should create task with valid data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		mockRepo.EXPECT().
			Create(match.PtrTo(&Task{
				Title:       "Test Task",
				Description: "Description",
				Status:      StatusNotStarted,
				Priority:    PriorityMedium,
			})).
			DoAndReturn(func(task *Task) error {
				task.ID = "TEST-ID"
				return nil
			})

		task, err := service.Create("Test Task", "Description", PriorityMedium, nil)

		assert.NoError(t, err)
		assert.Equal(t, "TEST-ID", task.ID)
		assert.Equal(t, "Test Task", task.Title)
		assert.Equal(t, "Description", task.Description)
	})

	t.Run("should return error with empty title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		// No expectations on mockRepo as it should not be called

		task, err := service.Create("", "Description", PriorityMedium, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty")
		assert.Nil(t, task)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)
		repoErr := errors.New("database error")

		mockRepo.EXPECT().
			Create(match.PtrTo(&Task{
				Title:       "Test Task",
				Description: "Description",
				Status:      StatusNotStarted,
				Priority:    PriorityMedium,
			})).
			Return(repoErr)

		task, err := service.Create("Test Task", "Description", PriorityMedium, nil)

		assert.Error(t, err)
		assert.Nil(t, task)
		assert.ErrorContains(t, err, "creating task")
	})
}

func TestService_Get(t *testing.T) {
	t.Run("should get task with valid ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)
		expectedTask := &Task{ID: "TEST-ID", Title: "Test Task"}

		mockRepo.EXPECT().
			Get("TEST-ID").
			Return(expectedTask, nil)

		task, err := service.Get("TEST-ID")

		assert.NoError(t, err)
		assert.Equal(t, expectedTask, task)
	})

	t.Run("should return error with empty ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		// No expectations on mockRepo as it should not be called

		task, err := service.Get("")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty")
		assert.Nil(t, task)
	})

	t.Run("should return repository error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		mockRepo.EXPECT().
			Get("TEST-ID").
			Return(nil, ErrTaskNotFound)

		task, err := service.Get("TEST-ID")

		assert.Error(t, err)
		assert.Nil(t, task)
	})
}

func TestService_List(t *testing.T) {
	t.Run("should return all tasks from repository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)
		tasks := []*Task{
			{ID: "1", Title: "Task 1"},
			{ID: "2", Title: "Task 2"},
		}

		mockRepo.EXPECT().
			List().
			Return(tasks, nil)

		result, err := service.List()

		assert.NoError(t, err)
		assert.Equal(t, tasks, result)
	})

	t.Run("should return repository error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)
		repoErr := errors.New("database error")

		mockRepo.EXPECT().
			List().
			Return(nil, repoErr)

		tasks, err := service.List()

		assert.Error(t, err)
		assert.Nil(t, tasks)
	})
}

func TestService_Update(t *testing.T) {
	t.Run("should update valid task", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)
		task := &Task{ID: "TEST-ID", Title: "Test Task"}

		mockRepo.EXPECT().
			Update(task).
			Return(nil)

		err := service.Update(task)

		assert.NoError(t, err)
	})

	t.Run("should call repository update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)
		task := &Task{
			ID:    "TEST-ID",
			Title: "Test Task",
		}

		mockRepo.EXPECT().
			Update(match.PtrTo(&Task{
				ID:    "TEST-ID",
				Title: "Test Task",
			})).
			Return(nil)

		err := service.Update(task)

		assert.NoError(t, err)
	})

	t.Run("should return error for nil task", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		// No expectations on mockRepo as it should not be called

		err := service.Update(nil)

		assert.ErrorIs(t, err, ErrInvalidTask)
	})

	t.Run("should return repository error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)
		task := &Task{ID: "TEST-ID", Title: "Test Task"}
		repoErr := errors.New("database error")

		mockRepo.EXPECT().
			Update(match.PtrTo(&Task{
				ID:    "TEST-ID",
				Title: "Test Task",
			})).
			Return(repoErr)

		err := service.Update(task)

		assert.Error(t, err)
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("should delete task with valid ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		mockRepo.EXPECT().
			Delete("TEST-ID").
			Return(nil)

		err := service.Delete("TEST-ID")

		assert.NoError(t, err)
	})

	t.Run("should return error with empty ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		// No expectations on mockRepo as it should not be called

		err := service.Delete("")

		assert.Error(t, err)
	})

	t.Run("should return repository error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)
		repoErr := errors.New("database error")

		mockRepo.EXPECT().
			Delete("TEST-ID").
			Return(repoErr)

		err := service.Delete("TEST-ID")

		assert.Error(t, err)
	})
}
