package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/utsabbera/task-master/core/chat"
	"github.com/utsabbera/task-master/core/task"
	"github.com/utsabbera/task-master/pkg/match"
	"github.com/utsabbera/task-master/pkg/util"
)

func TestHandler_Create(t *testing.T) {
	t.Run("should create task with valid data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		input := TaskInput{
			Title:       "Test Task",
			Description: "Test Description",
			Status:      task.StatusNotStarted,
			Priority:    util.Ptr(task.PriorityMedium),
			DueDate:     util.Ptr(time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)),
		}
		inputBytes, err := json.Marshal(input)
		require.NoError(t, err)

		createTime := time.Now().Truncate(time.Second)

		mockTaskService.EXPECT().Create(match.PtrTo(task.Task{
			Title:       input.Title,
			Description: input.Description,
			Status:      input.Status,
			Priority:    input.Priority,
			DueDate:     input.DueDate,
		})).DoAndReturn(func(tk *task.Task) error {
			tk.ID = "task-123"
			tk.CreatedAt = createTime
			tk.UpdatedAt = tk.CreatedAt
			return nil
		})

		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(inputBytes))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		handler.Create(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		var response Task
		err = json.Unmarshal(res.Body.Bytes(), &response)
		require.NoError(t, err)

		expected := Task{
			ID:          "task-123",
			Title:       input.Title,
			Description: input.Description,
			Status:      task.StatusNotStarted,
			Priority:    input.Priority,
			DueDate:     input.DueDate,
			CreatedAt:   createTime,
			UpdatedAt:   createTime,
		}
		assert.Equal(t, expected, response)
	})

	t.Run("should return bad request with invalid JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		handler.Create(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Invalid request payload")
	})

	t.Run("should return bad request when service returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		input := TaskInput{
			Title:       "Test Task",
			Description: "Test Description",
			Priority:    util.Ptr(task.PriorityMedium),
		}
		inputBytes, err := json.Marshal(input)
		require.NoError(t, err)

		mockTaskService.EXPECT().Create(match.PtrTo(task.Task{
			Title:       input.Title,
			Description: input.Description,
			Priority:    input.Priority,
		})).Return(fmt.Errorf("service error"))

		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(inputBytes))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		handler.Create(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "service error")
	})

	t.Run("should return bad request when required fields are missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		input := TaskInput{Description: "desc"}
		inputBytes, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(inputBytes))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		handler.Create(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Title cannot be empty")
	})
}

func TestHandler_Get(t *testing.T) {
	t.Run("should return task with valid ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		taskID := "task-123"
		existingTask := &task.Task{
			ID:          taskID,
			Title:       "Test Task",
			Description: "Test Description",
			Status:      task.StatusNotStarted,
		}
		mockTaskService.EXPECT().Get(taskID).Return(existingTask, nil)

		req := httptest.NewRequest(http.MethodGet, "/tasks/"+taskID, nil)
		req.SetPathValue("id", taskID)
		res := httptest.NewRecorder()
		handler.Get(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		var response Task
		err := json.Unmarshal(res.Body.Bytes(), &response)
		require.NoError(t, err)

		expected := Task{
			ID:          taskID,
			Title:       existingTask.Title,
			Description: existingTask.Description,
			Status:      task.StatusNotStarted,
		}
		assert.Equal(t, expected, response)
	})

	t.Run("should return bad request with missing ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		req := httptest.NewRequest(http.MethodGet, "/tasks/", nil)
		res := httptest.NewRecorder()
		handler.Get(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Missing task ID")
	})

	t.Run("should return not found when service returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		taskID := "non-existent"
		mockTaskService.EXPECT().Get(taskID).Return(nil, task.ErrTaskNotFound)

		req := httptest.NewRequest(http.MethodGet, "/tasks/"+taskID, nil)
		req.SetPathValue("id", taskID)
		res := httptest.NewRecorder()
		handler.Get(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Contains(t, res.Body.String(), "Task not found")
	})

	t.Run("should return internal server error when service returns unexpected error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		taskID := "task-err"
		mockTaskService.EXPECT().Get(taskID).Return(nil, fmt.Errorf("unexpected error"))

		req := httptest.NewRequest(http.MethodGet, "/tasks/"+taskID, nil)
		req.SetPathValue("id", taskID)
		res := httptest.NewRecorder()
		handler.Get(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Contains(t, res.Body.String(), "unexpected error")
	})
}

func TestHandler_List(t *testing.T) {
	t.Run("should return all tasks", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		tasks := []*task.Task{
			{
				ID:          "task-1",
				Title:       "Task 1",
				Description: "Description 1",
				Status:      task.StatusNotStarted,
			},
			{
				ID:          "task-2",
				Title:       "Task 2",
				Description: "Description 2",
				Status:      task.StatusInProgress,
			},
		}
		mockTaskService.EXPECT().List().Return(tasks, nil)

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		res := httptest.NewRecorder()
		handler.List(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		var response []Task
		err := json.Unmarshal(res.Body.Bytes(), &response)
		require.NoError(t, err)

		expected := []Task{
			{
				ID:          "task-1",
				Title:       "Task 1",
				Description: "Description 1",
				Status:      task.StatusNotStarted,
			},
			{
				ID:          "task-2",
				Title:       "Task 2",
				Description: "Description 2",
				Status:      task.StatusInProgress,
			},
		}
		assert.Equal(t, expected, response)
	})

	t.Run("should return server error when service fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		mockTaskService.EXPECT().List().Return(nil, fmt.Errorf("database error"))

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		res := httptest.NewRecorder()
		handler.List(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Contains(t, res.Body.String(), "Failed to retrieve tasks")
	})
}

func TestHandler_Update(t *testing.T) {
	t.Run("should update task with valid data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		taskID := "task-123"
		input := TaskInput{
			Title:       "Updated Task",
			Description: "Updated Description",
			Priority:    util.Ptr(task.PriorityHigh),
			DueDate:     util.Ptr(time.Date(2025, 2, 1, 12, 0, 0, 0, time.UTC)),
		}
		inputBytes, err := json.Marshal(input)
		require.NoError(t, err)

		updated := &task.Task{
			ID:          taskID,
			Title:       input.Title,
			Description: input.Description,
			Status:      task.StatusNotStarted,
			Priority:    input.Priority,
			DueDate:     input.DueDate,
		}
		mockTaskService.EXPECT().Update(taskID, match.PtrTo(task.Task{
			Title:       input.Title,
			Description: input.Description,
			Priority:    input.Priority,
			DueDate:     input.DueDate,
		})).Return(updated, nil)

		req := httptest.NewRequest(http.MethodPatch, "/tasks/"+taskID, bytes.NewReader(inputBytes))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", taskID)
		res := httptest.NewRecorder()
		handler.Update(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		var response Task
		err = json.Unmarshal(res.Body.Bytes(), &response)
		require.NoError(t, err)

		expected := Task{
			ID:          taskID,
			Title:       input.Title,
			Description: input.Description,
			Status:      task.StatusNotStarted,
			Priority:    input.Priority,
			DueDate:     input.DueDate,
		}
		assert.Equal(t, expected, response)
	})

	t.Run("should return bad request with missing ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		input := TaskInput{
			Title:       "Updated Task",
			Description: "Updated Description",
		}
		inputBytes, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, "/tasks/", bytes.NewReader(inputBytes))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		handler.Update(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Missing task ID")
	})

	t.Run("should return bad request with invalid JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		taskID := "task-123"
		req := httptest.NewRequest(http.MethodPatch, "/tasks/"+taskID, bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", taskID)
		res := httptest.NewRecorder()
		handler.Update(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Invalid request payload")
	})

	t.Run("should return not found when task doesn't exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		taskID := "non-existent"
		input := TaskInput{
			Title:       "Updated Task",
			Description: "Updated Description",
		}
		inputBytes, err := json.Marshal(input)
		require.NoError(t, err)

		mockTaskService.EXPECT().Update(taskID, match.PtrTo(task.Task{
			Title:       input.Title,
			Description: input.Description,
		})).Return(nil, task.ErrTaskNotFound)

		req := httptest.NewRequest(http.MethodPatch, "/tasks/"+taskID, bytes.NewReader(inputBytes))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", taskID)
		res := httptest.NewRecorder()
		handler.Update(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Contains(t, res.Body.String(), "Task not found")
	})

	t.Run("should return server error when update fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		taskID := "task-123"
		input := TaskInput{
			Title:       "Updated Task",
			Description: "Updated Description",
		}
		inputBytes, err := json.Marshal(input)
		require.NoError(t, err)

		mockTaskService.EXPECT().Update(taskID, match.PtrTo(task.Task{
			Title:       input.Title,
			Description: input.Description,
		})).Return(nil, fmt.Errorf("database error"))

		req := httptest.NewRequest(http.MethodPatch, "/tasks/"+taskID, bytes.NewReader(inputBytes))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", taskID)
		res := httptest.NewRecorder()
		handler.Update(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Contains(t, res.Body.String(), "database error")
	})

	t.Run("should update only provided fields (partial update)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		taskID := "task-123"
		input := TaskInput{
			Title: "Partial Update",
		}
		inputBytes, err := json.Marshal(input)
		require.NoError(t, err)

		updated := &task.Task{
			ID:    taskID,
			Title: input.Title,
			// other fields remain unchanged or zero
		}
		mockTaskService.EXPECT().Update(taskID, match.PtrTo(task.Task{
			Title: input.Title,
		})).Return(updated, nil)

		req := httptest.NewRequest(http.MethodPatch, "/tasks/"+taskID, bytes.NewReader(inputBytes))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", taskID)
		res := httptest.NewRecorder()
		handler.Update(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		var response Task
		err = json.Unmarshal(res.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, updated.ID, response.ID)
		assert.Equal(t, updated.Title, response.Title)
	})
}

func TestHandler_Delete(t *testing.T) {
	t.Run("should delete task with valid ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		taskID := "task-123"
		mockTaskService.EXPECT().Delete(taskID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/tasks/"+taskID, nil)
		req.SetPathValue("id", taskID)
		res := httptest.NewRecorder()
		handler.Delete(res, req)

		assert.Equal(t, http.StatusNoContent, res.Code)
		assert.Empty(t, res.Body.String())
	})

	t.Run("should return bad request with missing ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		req := httptest.NewRequest(http.MethodDelete, "/tasks/", nil)
		res := httptest.NewRecorder()
		handler.Delete(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Missing task ID")
	})

	t.Run("should return not found when task doesn't exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		taskID := "non-existent"
		mockTaskService.EXPECT().Delete(taskID).Return(task.ErrTaskNotFound)

		req := httptest.NewRequest(http.MethodDelete, "/tasks/"+taskID, nil)
		req.SetPathValue("id", taskID)
		res := httptest.NewRecorder()
		handler.Delete(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Contains(t, res.Body.String(), "Task not found")
	})

	t.Run("should return internal server error when service returns unexpected error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		taskID := "task-err"
		mockTaskService.EXPECT().Delete(taskID).Return(fmt.Errorf("unexpected error"))

		req := httptest.NewRequest(http.MethodDelete, "/tasks/"+taskID, nil)
		req.SetPathValue("id", taskID)
		res := httptest.NewRecorder()
		handler.Delete(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Contains(t, res.Body.String(), "unexpected error")
	})
}

func TestHandler_ProcessPrompt(t *testing.T) {
	t.Run("should process prompt successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		input := PromptInput{
			Text: "Create a new task to finish the report",
		}
		mockPromptService.EXPECT().ProcessPrompt(input.Text).Return("Task created: TASK-123", nil)

		body, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/prompts", bytes.NewReader(body))
		w := httptest.NewRecorder()
		handler.ProcessPrompt(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PromptResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response.Response, "TASK-123")
	})

	t.Run("should handle empty prompt", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		input := PromptInput{
			Text: "",
		}
		body, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/prompts", bytes.NewReader(body))
		w := httptest.NewRecorder()
		handler.ProcessPrompt(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "cannot be empty")
	})

	t.Run("should handle prompt service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTaskService := task.NewMockService(ctrl)
		mockPromptService := chat.NewMockService(ctrl)
		handler := NewHandler(mockTaskService, mockPromptService)

		input := PromptInput{
			Text: "Invalid prompt",
		}
		mockPromptService.EXPECT().ProcessPrompt(input.Text).Return("", errors.New("failed to process prompt"))

		body, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/prompts", bytes.NewReader(body))
		w := httptest.NewRecorder()
		handler.ProcessPrompt(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "failed to process prompt")
	})
}
