package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/utsabbera/task-master/core/assistant"

	taskcore "github.com/utsabbera/task-master/core/task"
)

//go:generate mockgen -destination=handler_mock.go -package=api . Handler

// Handler defines the interface for handling HTTP requests related to tasks.
type Handler interface {
	// Create creates a new task.
	Create(w http.ResponseWriter, r *http.Request)

	// Get retrieves a task by its ID.
	Get(w http.ResponseWriter, r *http.Request)

	// List lists all tasks.
	List(w http.ResponseWriter, r *http.Request)

	// Update updates an existing task by its ID.
	Update(w http.ResponseWriter, r *http.Request)

	// Delete deletes a task by its ID.
	Delete(w http.ResponseWriter, r *http.Request)

	// Chat handles natural language messages for task management.
	Chat(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	task      taskcore.Service
	assistant assistant.Service
}

// NewHandler returns a new instance of Handler for task operations.
func NewHandler(taskService taskcore.Service, assistantService assistant.Service) Handler {
	return &handler{
		task:      taskService,
		assistant: assistantService,
	}
}

// Create godoc
// @Summary Create Task
// @Description Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body TaskInput true "Task input"
// @Success 201 {object} Task
// @Router /tasks [post]
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var input TaskInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := r.Body.Close()
	if err != nil {
		http.Error(w, "Failed to close request body", http.StatusInternalServerError)
		return
	}

	if input.Title == "" {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}

	task := &taskcore.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
	}

	err = h.task.Create(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := mapTaskToResponse(task)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// Get godoc
// @Summary Get Task
// @Description Get a task by ID
// @Tags tasks
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} Task
// @Failure 404 {string} string "Task not found"
// @Router /tasks/{id} [get]
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	task, err := h.task.Get(id)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := mapTaskToResponse(task)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// List godoc
// @Summary List Tasks
// @Description List all tasks
// @Tags tasks
// @Produce json
// @Success 200 {array} Task
// @Router /tasks [get]
func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.task.List()
	if err != nil {
		http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := mapTasksToResponse(tasks)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// Update godoc
// @Summary Update Task
// @Description Partially update a task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body TaskInput true "Task fields to update"
// @Success 200 {object} Task
// @Failure 404 {string} string "Task not found"
// @Router /tasks/{id} [patch]
func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	var input TaskInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := r.Body.Close()
	if err != nil {
		http.Error(w, "Failed to close request body", http.StatusInternalServerError)
		return
	}

	patch := &taskcore.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
	}

	task, err := h.task.Update(id, patch)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := mapTaskToResponse(task)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// Delete godoc
// @Summary Delete Task
// @Description Delete a task by ID
// @Tags tasks
// @Param id path string true "Task ID"
// @Success 204 {string} string "Task deleted"
// @Failure 404 {string} string "Task not found"
// @Router /tasks/{id} [delete]
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	if err := h.task.Delete(id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Chat godoc
// @Summary Chat
// @Description Chat in natural language for task management
// @Tags chat
// @Accept json
// @Produce json
// @Param chat body ChatInput true "Chat input"
// @Success 200 {object} ChatResponse
// @Router /chat [post]
func (h *handler) Chat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input ChatInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Text == "" {
		http.Error(w, "Chat text cannot be empty", http.StatusBadRequest)
		return
	}

	response, err := h.assistant.Chat(ctx, input.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := ChatResponse{
		Response: response,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, taskcore.ErrTaskNotFound) {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}
