package api

import (
	"time"

	"github.com/utsabbera/task-master/core/task"
)

// Task represents a task in the task management system.
type Task struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Status      task.Status    `json:"status"`
	Priority    *task.Priority `json:"priority"`
	DueDate     *time.Time     `json:"dueDate"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

type TaskInput struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Status      task.Status    `json:"status"`
	Priority    *task.Priority `json:"priority"`
	DueDate     *time.Time     `json:"dueDate"`
}

// ChatInput represents a natural language message for task management.
type ChatInput struct {
	Text string `json:"text"`
}

// ChatResponse represents the response to a natural language message.
type ChatResponse struct {
	Response string `json:"response"`
}
