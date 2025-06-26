package assistant

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utsabbera/task-master/core/task"
	"github.com/utsabbera/task-master/pkg/assistant"
	"go.uber.org/mock/gomock"
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskService := task.NewMockService(ctrl)
	mockAssistant := assistant.NewMockClient(ctrl)

	mockAssistant.EXPECT().Init()

	svc := NewService(mockTaskService, mockAssistant)
	assert.NotNil(t, svc)
}

func TestService_Chat(t *testing.T) {
	t.Run("should create task from message", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTaskService := task.NewMockService(ctrl)
		mockAssistant := assistant.NewMockClient(ctrl)

		mockAssistant.EXPECT().Init()

		service := NewService(mockTaskService, mockAssistant)
		message := "Create a new task"

		expectedTask := &task.Task{
			Title:       message,
			Description: "",
			Priority:    nil,
			DueDate:     nil,
		}
		mockTaskService.EXPECT().Create(expectedTask).DoAndReturn(func(t *task.Task) error {
			t.ID = "TASK-123"
			return nil
		})

		result, err := service.Chat(ctx, message)
		assert.NoError(t, err)
		assert.Contains(t, result, "Task created")
	})

	t.Run("should handle errors from task service", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTaskService := task.NewMockService(ctrl)
		mockAssistant := assistant.NewMockClient(ctrl)

		mockAssistant.EXPECT().Init()

		service := NewService(mockTaskService, mockAssistant)
		message := "Invalid task"

		expectedTask := &task.Task{
			Title:       message,
			Description: "",
			Priority:    nil,
			DueDate:     nil,
		}

		mockTaskService.EXPECT().Create(expectedTask).Return(errors.New("task creation failed"))

		result, err := service.Chat(ctx, message)
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.ErrorContains(t, err, "task creation failed")
	})

	// TODO: Fix the testcases
}
