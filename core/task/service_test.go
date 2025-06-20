package task

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/utsabbera/task-master/pkg/match"
	"go.uber.org/mock/gomock"
)

func TestTaskService_Create(t *testing.T) {
	t.Run("should create task with valid data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		priority := PriorityMedium
		due := time.Now().Add(24 * time.Hour).Truncate(time.Second)

		mockRepo.EXPECT().
			Create(match.PtrTo(&Task{
				Title:       "Test Task",
				Description: "Description",
				Status:      StatusNotStarted,
				Priority:    &priority,
				DueDate:     &due,
			})).
			DoAndReturn(func(task *Task) error {
				task.ID = "TEST-ID"
				return nil
			})

		task, err := service.Create("Test Task", "Description", &priority, &due)

		assert.NoError(t, err)
		assert.Equal(t, "TEST-ID", task.ID)
		assert.Equal(t, "Test Task", task.Title)
		assert.Equal(t, "Description", task.Description)
		assert.Equal(t, &priority, task.Priority)
		assert.WithinDuration(t, due, *task.DueDate, time.Second)
	})

	t.Run("should create task without priority", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		due := time.Now().Add(24 * time.Hour).Truncate(time.Second)

		mockRepo.EXPECT().
			Create(match.PtrTo(&Task{
				Title:       "No Priority",
				Description: "Description",
				Status:      StatusNotStarted,
				Priority:    nil,
				DueDate:     &due,
			})).
			DoAndReturn(func(task *Task) error {
				task.ID = "TEST-ID"
				return nil
			})

		task, err := service.Create("No Priority", "Description", nil, &due)

		assert.NoError(t, err)
		assert.Equal(t, "TEST-ID", task.ID)
		assert.Nil(t, task.Priority)
	})

	t.Run("should create task without due date", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		priority := PriorityMedium

		mockRepo.EXPECT().
			Create(match.PtrTo(&Task{
				Title:       "No DueDate",
				Description: "Description",
				Status:      StatusNotStarted,
				Priority:    &priority,
				DueDate:     nil,
			})).
			DoAndReturn(func(task *Task) error {
				task.ID = "TEST-ID"
				return nil
			})

		task, err := service.Create("No DueDate", "Description", &priority, nil)

		assert.NoError(t, err)
		assert.Equal(t, "TEST-ID", task.ID)
		assert.Nil(t, task.DueDate)
	})

	t.Run("should return error for empty title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		task, err := service.Create("", "Description", nil, nil)

		assert.Error(t, err)
		assert.Nil(t, task)
		assert.Contains(t, err.Error(), "title cannot be empty")
	})

	t.Run("should return error when repository create fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		mockRepo.EXPECT().
			Create(gomock.Any()).
			Return(errors.New("repository error"))

		task, err := service.Create("Test Task", "Description", nil, nil)

		assert.Error(t, err)
		assert.Nil(t, task)
		assert.Contains(t, err.Error(), "repository error")
	})
}

func TestTaskService_Get(t *testing.T) {
	t.Run("should get task by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		mockRepo.EXPECT().
			Get("TEST-ID").
			Return(&Task{ID: "TEST-ID", Title: "Test Task"}, nil)

		task, err := service.Get("TEST-ID")

		assert.NoError(t, err)
		assert.Equal(t, "TEST-ID", task.ID)
		assert.Equal(t, "Test Task", task.Title)
	})

	t.Run("should return error when task not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		mockRepo.EXPECT().
			Get("UNKNOWN").
			Return(nil, errors.New("not found"))

		task, err := service.Get("UNKNOWN")

		assert.Error(t, err)
		assert.Nil(t, task)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestTaskService_List(t *testing.T) {
	t.Run("should list all tasks", func(t *testing.T) {
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
		assert.Len(t, result, 2)
	})

	t.Run("should return empty list when no tasks", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		mockRepo.EXPECT().
			List().
			Return([]*Task{}, nil)

		result, err := service.List()

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		mockRepo.EXPECT().
			List().
			Return(nil, errors.New("repository error"))

		result, err := service.List()

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "repository error")
	})
}

func TestTaskService_Update(t *testing.T) {
	t.Run("should update task", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		task := &Task{ID: "TEST-ID", Title: "Updated Task"}

		mockRepo.EXPECT().
			Update(task).
			Return(nil)

		err := service.Update(task)

		assert.NoError(t, err)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		task := &Task{ID: "TEST-ID", Title: "Updated Task"}

		mockRepo.EXPECT().
			Update(task).
			Return(errors.New("update error"))

		err := service.Update(task)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update error")
	})
}

func TestTaskService_Delete(t *testing.T) {
	t.Run("should delete task", func(t *testing.T) {
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

	t.Run("should return error when delete fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		mockRepo.EXPECT().
			Delete("TEST-ID").
			Return(errors.New("delete error"))

		err := service.Delete("TEST-ID")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete error")
	})
}
