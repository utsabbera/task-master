package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/utsabbera/task-master/core"
	"github.com/utsabbera/task-master/pkg/match"
	"github.com/utsabbera/task-master/pkg/util"
	"go.uber.org/mock/gomock"
)

func TestHandler_Create(t *testing.T) {
	t.Run("should create task with valid data", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := core.NewMockTaskService(ctrl)
		handler := NewHandler(mockService)

		input := TaskInput{
			Title:       "Test Task",
			Description: "Test Description",
			Priority:    util.Ptr(core.PriorityMedium),
			DueDate:     util.Ptr(time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)),
		}

		inputBytes, err := json.Marshal(input)
		require.NoError(t, err)
		mockService.EXPECT().
			Create(
				input.Title,
				input.Description,
				input.Priority,
				input.DueDate,
			).
			Return(&core.Task{
				ID:          "task-123",
				Title:       input.Title,
				Description: input.Description,
				Status:      core.StatusNotStarted,
				Priority:    input.Priority,
				DueDate:     input.DueDate,
			}, nil)

		// Act
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(inputBytes))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.Create(res, req)

		// Assert
		assert.Equal(t, http.StatusCreated, res.Code)

		var response Task
		err = json.Unmarshal(res.Body.Bytes(), &response)
		require.NoError(t, err)

		expected := Task{
			ID:          "task-123",
			Title:       input.Title,
			Description: input.Description,
			Status:      core.StatusNotStarted,
			Priority:    input.Priority,
			DueDate:     input.DueDate,
		}
		assert.Equal(t, expected, response)
	})

	t.Run("should return bad request with invalid JSON", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := core.NewMockTaskService(ctrl)
		handler := NewHandler(mockService)

		// Act
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.Create(res, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Invalid request payload")
	})

	t.Run("should return bad request when service returns error", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := core.NewMockTaskService(ctrl)
		handler := NewHandler(mockService)

		input := TaskInput{
			Title:       "Test Task",
			Description: "Test Description",
			Priority:    util.Ptr(core.PriorityMedium),
		}

		inputBytes, err := json.Marshal(input)
		require.NoError(t, err)

		mockService.EXPECT().
			Create(
				input.Title,
				input.Description,
				input.Priority,
				nil,
			).
			Return(nil, fmt.Errorf("service error"))

		// Act
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(inputBytes))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.Create(res, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "service error")
	})
}

func TestHandler_Get(t *testing.T) {
	t.Run("should return task with valid ID", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := core.NewMockTaskService(ctrl)
		handler := NewHandler(mockService)

		taskID := "task-123"
		task := &core.Task{
			ID:          taskID,
			Title:       "Test Task",
			Description: "Test Description",
			Status:      core.StatusNotStarted,
		}

		mockService.EXPECT().
			Get(taskID).
			Return(task, nil)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/tasks/"+taskID, nil)
		req.SetPathValue("id", taskID)

		res := httptest.NewRecorder()

		handler.Get(res, req)

		// Assert
		assert.Equal(t, http.StatusOK, res.Code)

		var response Task
		err := json.Unmarshal(res.Body.Bytes(), &response)
		require.NoError(t, err)

		expected := Task{
			ID:          taskID,
			Title:       task.Title,
			Description: task.Description,
			Status:      core.StatusNotStarted,
		}
		assert.Equal(t, expected, response)
	})

	t.Run("should return bad request with missing ID", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := core.NewMockTaskService(ctrl)
		handler := NewHandler(mockService)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/tasks/", nil)
		res := httptest.NewRecorder()

		handler.Get(res, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Missing task ID")
	})

	t.Run("should return not found when service returns error", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := core.NewMockTaskService(ctrl)
		handler := NewHandler(mockService)

		taskID := "non-existent"

		mockService.EXPECT().
			Get(taskID).
			Return(nil, core.ErrTaskNotFound)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/tasks/"+taskID, nil)
		req.SetPathValue("id", taskID)

		res := httptest.NewRecorder()

		handler.Get(res, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Contains(t, res.Body.String(), "Task not found")
	})
}

func TestHandler_List(t *testing.T) {
	t.Run("should return all tasks", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := core.NewMockTaskService(ctrl)
		handler := NewHandler(mockService)

		tasks := []*core.Task{
			{
				ID:          "task-1",
				Title:       "Task 1",
				Description: "Description 1",
				Status:      core.StatusNotStarted,
			},
			{
				ID:          "task-2",
				Title:       "Task 2",
				Description: "Description 2",
				Status:      core.StatusInProgress,
			},
		}

		mockService.EXPECT().
			List().
			Return(tasks, nil)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		res := httptest.NewRecorder()

		handler.List(res, req)

		// Assert
		assert.Equal(t, http.StatusOK, res.Code)

		var response []Task
		err := json.Unmarshal(res.Body.Bytes(), &response)
		require.NoError(t, err)

		expected := []Task{
			{
				ID:          "task-1",
				Title:       "Task 1",
				Description: "Description 1",
				Status:      core.StatusNotStarted,
			},
			{
				ID:          "task-2",
				Title:       "Task 2",
				Description: "Description 2",
				Status:      core.StatusInProgress,
			},
		}
		assert.Equal(t, expected, response)
	})

	t.Run("should return server error when service fails", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := core.NewMockTaskService(ctrl)
		handler := NewHandler(mockService)

		mockService.EXPECT().
			List().
			Return(nil, fmt.Errorf("database error"))

		// Act
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		res := httptest.NewRecorder()

		handler.List(res, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Contains(t, res.Body.String(), "Failed to retrieve tasks")
	})
}

func TestHandler_Update(t *testing.T) {
	t.Run("should update task with valid data", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := core.NewMockTaskService(ctrl)
		handler := NewHandler(mockService)

		taskID := "task-123"
		input := TaskInput{
			Title:       "Updated Task",
			Description: "Updated Description",
			Priority:    util.Ptr(core.PriorityHigh),
			DueDate:     util.Ptr(time.Date(2025, 2, 1, 12, 0, 0, 0, time.UTC)),
		}

		inputBytes, err := json.Marshal(input)
		require.NoError(t, err)

		mockService.EXPECT().
			Get(taskID).
			Return(&core.Task{
				ID:          taskID,
				Title:       "Original Task",
				Description: "Original Description",
				Status:      core.StatusNotStarted,
			}, nil)

		mockService.EXPECT().
			Update(match.PtrTo(core.Task{
				ID:          taskID,
				Title:       input.Title,
				Description: input.Description,
				Status:      core.StatusNotStarted,
				Priority:    input.Priority,
				DueDate:     input.DueDate,
			})).
			Return(nil)

		// Act
		req := httptest.NewRequest(http.MethodPut, "/tasks/"+taskID, bytes.NewReader(inputBytes))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", taskID)
		res := httptest.NewRecorder()

		handler.Update(res, req)

		// Assert
		assert.Equal(t, http.StatusOK, res.Code)

		var response Task
		err = json.Unmarshal(res.Body.Bytes(), &response)
		require.NoError(t, err)

		expected := Task{
			ID:          taskID,
			Title:       input.Title,
			Description: input.Description,
			Status:      core.StatusNotStarted,
			Priority:    input.Priority,
			DueDate:     input.DueDate,
		}
		assert.Equal(t, expected, response)
	})

	t.Run("should return bad request with missing ID", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := core.NewMockTaskService(ctrl)
		handler := NewHandler(mockService)

		input := TaskInput{
			Title:       "Updated Task",
			Description: "Updated Description",
		}

		inputBytes, err := json.Marshal(input)
		require.NoError(t, err)

		// Act
		req := httptest.NewRequest(http.MethodPut, "/tasks/", bytes.NewReader(inputBytes))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.Update(res, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Missing task ID")
	})

	t.Run("should return bad request with invalid JSON", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := core.NewMockTaskService(ctrl)
		handler := NewHandler(mockService)

		taskID := "task-123"

		// Act
		req := httptest.NewRequest(http.MethodPut, "/tasks/"+taskID, bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", taskID)
		res := httptest.NewRecorder()

		handler.Update(res, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Invalid request payload")
	})

	t.Run("should return not found when task doesn't exist", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := core.NewMockTaskService(ctrl)
		handler := NewHandler(mockService)

		taskID := "non-existent"
		input := TaskInput{
			Title:       "Updated Task",
			Description: "Updated Description",
		}

		inputBytes, err := json.Marshal(input)
		require.NoError(t, err)

		mockService.EXPECT().
			Get(taskID).
			Return(nil, core.ErrTaskNotFound)

		// Act
		req := httptest.NewRequest(http.MethodPut, "/tasks/"+taskID, bytes.NewReader(inputBytes))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", taskID)
		res := httptest.NewRecorder()

		handler.Update(res, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Contains(t, res.Body.String(), "Task not found")
	})

	t.Run("should return server error when update fails", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := core.NewMockTaskService(ctrl)
		handler := NewHandler(mockService)

		taskID := "task-123"
		input := TaskInput{
			Title:       "Updated Task",
			Description: "Updated Description",
		}

		inputBytes, err := json.Marshal(input)
		require.NoError(t, err)

		mockService.EXPECT().
			Get(taskID).
			Return(&core.Task{
				ID:          taskID,
				Title:       "Original Task",
				Description: "Original Description",
				Status:      core.StatusNotStarted,
			}, nil)

		mockService.EXPECT().
			Update(gomock.Any()).
			Return(fmt.Errorf("database error"))

		// Act
		req := httptest.NewRequest(http.MethodPut, "/tasks/"+taskID, bytes.NewReader(inputBytes))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", taskID)
		res := httptest.NewRecorder()

		handler.Update(res, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Contains(t, res.Body.String(), "database error")
	})
}

func TestHandler_Delete(t *testing.T) {
	t.Run("should delete task with valid ID", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := core.NewMockTaskService(ctrl)
		handler := NewHandler(mockService)

		taskID := "task-123"

		mockService.EXPECT().
			Delete(taskID).
			Return(nil)

		// Act
		req := httptest.NewRequest(http.MethodDelete, "/tasks/"+taskID, nil)
		req.SetPathValue("id", taskID)

		res := httptest.NewRecorder()

		handler.Delete(res, req)

		// Assert
		assert.Equal(t, http.StatusNoContent, res.Code)
		assert.Empty(t, res.Body.String())
	})

	t.Run("should return bad request with missing ID", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := core.NewMockTaskService(ctrl)
		handler := NewHandler(mockService)

		// Act
		req := httptest.NewRequest(http.MethodDelete, "/tasks/", nil)
		res := httptest.NewRecorder()

		handler.Delete(res, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Missing task ID")
	})

	t.Run("should return not found when task doesn't exist", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := core.NewMockTaskService(ctrl)
		handler := NewHandler(mockService)

		taskID := "non-existent"

		mockService.EXPECT().
			Delete(taskID).
			Return(core.ErrTaskNotFound)

		// Act
		req := httptest.NewRequest(http.MethodDelete, "/tasks/"+taskID, nil)
		req.SetPathValue("id", taskID)

		res := httptest.NewRecorder()

		handler.Delete(res, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Contains(t, res.Body.String(), "Task not found")
	})
}
