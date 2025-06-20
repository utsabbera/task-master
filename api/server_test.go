package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/utsabbera/task-master/core"
	"github.com/utsabbera/task-master/pkg/util"
)

func TestNewServer(t *testing.T) {
	t.Run("should return configured HTTP server", func(t *testing.T) {
		cfg := ServerConfig{Addr: ":9090"}
		server := NewServer(cfg)

		assert.NotNil(t, server)
		assert.Equal(t, ":9090", server.Addr)
		assert.NotNil(t, server.Handler)
	})

	t.Run("should use default Addr when empty", func(t *testing.T) {
		cfg := ServerConfig{Addr: ""}
		server := NewServer(cfg)

		assert.Equal(t, ":8080", server.Addr)
	})
}

func TestIntegration_Server(t *testing.T) {
	t.Run("should create task with title and description", func(t *testing.T) {
		repo := core.NewDefaultMemoryRepository()
		service := core.NewTaskService(repo)
		promptService := core.NewPromptService(service)
		handler := NewHandler(service, promptService)
		router := NewRouter(handler)

		ts := httptest.NewServer(router)
		defer ts.Close()

		input := TaskInput{
			Title:       "Test Task",
			Description: "Test description",
		}
		body, err := json.Marshal(input)
		require.NoError(t, err)

		resp, err := http.Post(ts.URL+"/tasks", "application/json", bytes.NewReader(body))
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())

		var task Task
		err = json.Unmarshal(respBody, &task)
		require.NoError(t, err)

		assert.Equal(t, input.Title, task.Title)
		assert.Equal(t, input.Description, task.Description)
		assert.NotEmpty(t, task.ID)
	})

	t.Run("should create task with due date", func(t *testing.T) {
		repo := core.NewDefaultMemoryRepository()
		service := core.NewTaskService(repo)
		promptService := core.NewPromptService(service)
		handler := NewHandler(service, promptService)
		router := NewRouter(handler)

		ts := httptest.NewServer(router)
		defer ts.Close()

		input := TaskInput{
			Title:       "Test Task",
			Description: "Test description",
			DueDate:     util.Ptr(time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)),
		}
		body, err := json.Marshal(input)
		require.NoError(t, err)

		resp, err := http.Post(ts.URL+"/tasks", "application/json", bytes.NewReader(body))
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())

		var task Task
		err = json.Unmarshal(respBody, &task)
		require.NoError(t, err)

		assert.Equal(t, input.Title, task.Title)
		assert.Equal(t, input.Description, task.Description)
		assert.Equal(t, input.DueDate, task.DueDate)
		assert.NotEmpty(t, task.ID)
	})

	t.Run("should create task with priority", func(t *testing.T) {
		repo := core.NewDefaultMemoryRepository()
		service := core.NewTaskService(repo)
		promptService := core.NewPromptService(service)
		handler := NewHandler(service, promptService)
		router := NewRouter(handler)

		ts := httptest.NewServer(router)
		defer ts.Close()

		input := TaskInput{
			Title:       "Test Task",
			Description: "Test description",
			Priority:    util.Ptr(core.PriorityHigh),
		}
		body, err := json.Marshal(input)
		require.NoError(t, err)

		resp, err := http.Post(ts.URL+"/tasks", "application/json", bytes.NewReader(body))
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())

		var task Task
		err = json.Unmarshal(respBody, &task)
		require.NoError(t, err)

		assert.Equal(t, input.Title, task.Title)
		assert.Equal(t, input.Description, task.Description)
		assert.Equal(t, input.Priority, task.Priority)
		assert.NotEmpty(t, task.ID)
	})

	t.Run("should get created task with Id", func(t *testing.T) {
		repo := core.NewDefaultMemoryRepository()
		service := core.NewTaskService(repo)
		promptService := core.NewPromptService(service)
		handler := NewHandler(service, promptService)
		router := NewRouter(handler)

		ts := httptest.NewServer(router)
		defer ts.Close()

		input := TaskInput{
			Title:       "Test Task",
			Description: "Test description",
		}
		body, err := json.Marshal(input)
		require.NoError(t, err)

		resp, err := http.Post(ts.URL+"/tasks", "application/json", bytes.NewReader(body))
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())

		var created Task
		err = json.Unmarshal(respBody, &created)
		require.NoError(t, err)
		require.NotEmpty(t, created.ID)

		resp, err = http.Get(ts.URL + "/tasks/" + created.ID)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		getBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())

		var fetched Task
		err = json.Unmarshal(getBody, &fetched)
		require.NoError(t, err)

		assert.Equal(t, created.ID, fetched.ID)
		assert.Equal(t, input.Title, fetched.Title)
		assert.Equal(t, input.Description, fetched.Description)
	})

	t.Run("should get list of created tasks", func(t *testing.T) {
		repo := core.NewDefaultMemoryRepository()
		service := core.NewTaskService(repo)
		promptService := core.NewPromptService(service)
		handler := NewHandler(service, promptService)
		router := NewRouter(handler)

		ts := httptest.NewServer(router)
		defer ts.Close()

		inputs := []TaskInput{
			{Title: "Task 1", Description: "Desc 1"},
			{Title: "Task 2", Description: "Desc 2"},
		}
		var createdIDs []string
		for _, input := range inputs {
			body, err := json.Marshal(input)
			require.NoError(t, err)

			resp, err := http.Post(ts.URL+"/tasks", "application/json", bytes.NewReader(body))
			require.NoError(t, err)

			assert.Equal(t, http.StatusCreated, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.NoError(t, resp.Body.Close())

			var task Task
			err = json.Unmarshal(respBody, &task)
			require.NoError(t, err)
			createdIDs = append(createdIDs, task.ID)
		}

		resp, err := http.Get(ts.URL + "/tasks")
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())

		var tasks []Task
		err = json.Unmarshal(respBody, &tasks)
		require.NoError(t, err)

		assert.Len(t, tasks, len(inputs))
		gotIDs := make([]string, len(tasks))
		for i, task := range tasks {
			gotIDs[i] = task.ID
		}
		assert.ElementsMatch(t, createdIDs, gotIDs)
	})

	t.Run("should update created task", func(t *testing.T) {
		repo := core.NewDefaultMemoryRepository()
		taskService := core.NewTaskService(repo)
		promptService := core.NewPromptService(taskService)
		handler := NewHandler(taskService, promptService)
		router := NewRouter(handler)

		ts := httptest.NewServer(router)
		defer ts.Close()

		input := TaskInput{
			Title:       "Initial Title",
			Description: "Initial Description",
		}
		body, err := json.Marshal(input)
		require.NoError(t, err)

		createResp, err := http.Post(ts.URL+"/tasks", "application/json", bytes.NewReader(body))
		require.NoError(t, err)

		require.Equal(t, http.StatusCreated, createResp.StatusCode)
		createRespBody, err := io.ReadAll(createResp.Body)
		require.NoError(t, err)
		require.NoError(t, createResp.Body.Close())

		var created Task
		err = json.Unmarshal(createRespBody, &created)
		require.NoError(t, err)
		require.NotEmpty(t, created.ID)

		updateInput := TaskInput{
			Title:       "Updated Title",
			Description: "Updated Description",
			Priority:    util.Ptr(core.PriorityMedium),
		}
		updateBody, err := json.Marshal(updateInput)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, ts.URL+"/tasks/"+created.ID, bytes.NewReader(updateBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		updateResp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, updateResp.StatusCode)
		updatedRespBody, err := io.ReadAll(updateResp.Body)
		require.NoError(t, err)
		require.NoError(t, updateResp.Body.Close())

		var updated Task
		err = json.Unmarshal(updatedRespBody, &updated)
		require.NoError(t, err)

		assert.Equal(t, created.ID, updated.ID)
		assert.Equal(t, updateInput.Title, updated.Title)
		assert.Equal(t, updateInput.Description, updated.Description)
		assert.Equal(t, updateInput.Priority, updated.Priority)
	})

	t.Run("should delete created task", func(t *testing.T) {
		repo := core.NewDefaultMemoryRepository()
		service := core.NewTaskService(repo)
		promptService := core.NewPromptService(service)
		handler := NewHandler(service, promptService)
		router := NewRouter(handler)

		ts := httptest.NewServer(router)
		defer ts.Close()

		input := TaskInput{
			Title:       "Task to Delete",
			Description: "Delete me",
		}
		body, err := json.Marshal(input)
		require.NoError(t, err)

		createResp, err := http.Post(ts.URL+"/tasks", "application/json", bytes.NewReader(body))
		require.NoError(t, err)

		require.Equal(t, http.StatusCreated, createResp.StatusCode)
		createRespBody, err := io.ReadAll(createResp.Body)
		require.NoError(t, err)
		require.NoError(t, createResp.Body.Close())

		var created Task
		err = json.Unmarshal(createRespBody, &created)
		require.NoError(t, err)
		require.NotEmpty(t, created.ID)

		req, err := http.NewRequest(http.MethodDelete, ts.URL+"/tasks/"+created.ID, nil)
		require.NoError(t, err)

		delResp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, delResp.StatusCode)
		require.NoError(t, delResp.Body.Close())

		getResp, err := http.Get(ts.URL + "/tasks/" + created.ID)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
		require.NoError(t, getResp.Body.Close())
	})

	// TODO: Add tests for prompt processing
}
