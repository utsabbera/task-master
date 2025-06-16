package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utsabbera/task-master/pkg/match"
	"go.uber.org/mock/gomock"
)

func TestCreateTask(t *testing.T) {
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

		task, err := service.CreateTask("Test Task", "Description", PriorityMedium, nil)

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

		task, err := service.CreateTask("", "Description", PriorityMedium, nil)

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

		task, err := service.CreateTask("Test Task", "Description", PriorityMedium, nil)

		assert.Error(t, err)
		assert.Nil(t, task)
		assert.ErrorContains(t, err, "creating task")
	})
}

func TestGetTask(t *testing.T) {
	t.Run("should get task with valid ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)
		expectedTask := &Task{ID: "TEST-ID", Title: "Test Task"}

		mockRepo.EXPECT().
			FindByID("TEST-ID").
			Return(expectedTask, nil)

		task, err := service.GetTask("TEST-ID")

		assert.NoError(t, err)
		assert.Equal(t, expectedTask, task)
	})

	t.Run("should return error with empty ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		// No expectations on mockRepo as it should not be called

		task, err := service.GetTask("")

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
			FindByID("TEST-ID").
			Return(nil, ErrTaskNotFound)

		task, err := service.GetTask("TEST-ID")

		assert.Error(t, err)
		assert.Nil(t, task)
	})
}

func TestGetAllTasks(t *testing.T) {
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
			FindAll().
			Return(tasks, nil)

		result, err := service.GetAllTasks()

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
			FindAll().
			Return(nil, repoErr)

		tasks, err := service.GetAllTasks()

		assert.Error(t, err)
		assert.Nil(t, tasks)
	})
}

func TestUpdateTask(t *testing.T) {
	t.Run("should update valid task", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)
		task := &Task{ID: "TEST-ID", Title: "Test Task"}

		mockRepo.EXPECT().
			Update(task).
			Return(nil)

		err := service.UpdateTask(task)

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

		err := service.UpdateTask(task)

		assert.NoError(t, err)
	})

	t.Run("should return error for nil task", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		// No expectations on mockRepo as it should not be called

		err := service.UpdateTask(nil)

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

		err := service.UpdateTask(task)

		assert.Error(t, err)
	})
}

func TestDeleteTask(t *testing.T) {
	t.Run("should delete task with valid ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		mockRepo.EXPECT().
			Delete("TEST-ID").
			Return(nil)

		err := service.DeleteTask("TEST-ID")

		assert.NoError(t, err)
	})

	t.Run("should return error with empty ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)

		// No expectations on mockRepo as it should not be called

		err := service.DeleteTask("")

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

		err := service.DeleteTask("TEST-ID")

		assert.Error(t, err)
	})
}

func TestUpdateTaskStatus(t *testing.T) {
	t.Run("should update task status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)
		task := &Task{ID: "TEST-ID", Status: StatusNotStarted}

		mockRepo.EXPECT().
			FindByID("TEST-ID").
			Return(task, nil)

		mockRepo.EXPECT().
			Update(match.PtrTo(&Task{
				ID:     "TEST-ID",
				Status: StatusCompleted,
			})).
			Return(nil)

		err := service.UpdateTaskStatus("TEST-ID", StatusCompleted)

		assert.NoError(t, err)
		assert.Equal(t, StatusCompleted, task.Status)
	})

	t.Run("should call repository update after status change", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)
		task := &Task{
			ID:     "TEST-ID",
			Status: StatusNotStarted,
		}

		mockRepo.EXPECT().
			FindByID("TEST-ID").
			Return(task, nil)

		mockRepo.EXPECT().
			Update(match.PtrTo(&Task{
				ID:     "TEST-ID",
				Status: StatusInProgress,
			})).
			Return(nil)

		err := service.UpdateTaskStatus("TEST-ID", StatusInProgress)

		assert.NoError(t, err)
		assert.Equal(t, StatusInProgress, task.Status)
	})

	t.Run("should return error when finding task fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)
		repoErr := errors.New("database error")

		mockRepo.EXPECT().
			FindByID("TEST-ID").
			Return(nil, repoErr)

		// No Update expectations as it should not be called

		err := service.UpdateTaskStatus("TEST-ID", StatusCompleted)

		assert.Error(t, err)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := NewService(mockRepo)
		task := &Task{ID: "TEST-ID", Status: StatusNotStarted}
		repoErr := errors.New("database error")

		mockRepo.EXPECT().
			FindByID("TEST-ID").
			Return(task, nil)

		mockRepo.EXPECT().
			Update(gomock.Any()).
			Return(repoErr)

		err := service.UpdateTaskStatus("TEST-ID", StatusCompleted)

		assert.Error(t, err)
	})
}
