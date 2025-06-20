package chat

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utsabbera/task-master/core/task"
	"go.uber.org/mock/gomock"
)

func TestService_ProcessPrompt(t *testing.T) {
	t.Run("should create task from prompt", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTaskService := task.NewMockService(ctrl)
		promptService := NewService(mockTaskService)

		mockTaskService.EXPECT().
			Create("Create a new task", "", nil, nil).
			Return(&task.Task{
				ID:    "TASK-123",
				Title: "Create a new task",
			}, nil)

		result, err := promptService.ProcessPrompt("Create a new task")

		assert.NoError(t, err)
		assert.Contains(t, result, "TASK-123")
	})

	t.Run("should handle errors from task service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTaskService := task.NewMockService(ctrl)
		promptService := NewService(mockTaskService)

		mockTaskService.EXPECT().
			Create("Invalid task", "", nil, nil).
			Return(nil, errors.New("task creation failed"))

		result, err := promptService.ProcessPrompt("Invalid task")

		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Contains(t, err.Error(), "task creation failed")
	})
}
