package api

import (
	"net/http"
)

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body TaskInput true "Task input"
// @Success 201 {object} Task
// @Router /tasks [post]
func (h *taskHandler) Create(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// GetTask godoc
// @Summary Get a task by ID
// @Description Get a task by ID
// @Tags tasks
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} Task
// @Failure 404 {string} string "Task not found"
// @Router /tasks/{id} [get]
func (h *taskHandler) Get(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// ListTasks godoc
// @Summary List all tasks
// @Description List all tasks
// @Tags tasks
// @Produce json
// @Success 200 {array} Task
// @Router /tasks [get]
func (h *taskHandler) List(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// UpdateTask godoc
// @Summary Update a task by ID
// @Description Update a task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body TaskInput true "Task input"
// @Success 200 {object} Task
// @Failure 404 {string} string "Task not found"
// @Router /tasks/{id} [put]
func (h *taskHandler) Update(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// DeleteTask godoc
// @Summary Delete a task by ID
// @Description Delete a task by ID
// @Tags tasks
// @Param id path string true "Task ID"
// @Success 204 {string} string "Task deleted"
// @Failure 404 {string} string "Task not found"
// @Router /tasks/{id} [delete]
func (h *taskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type taskHandler struct {
}

func NewHandler() Handler {
	return &taskHandler{}
}
