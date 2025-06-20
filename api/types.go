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
	Priority    *task.Priority `json:"priority"`
	DueDate     *time.Time     `json:"dueDate"`
}

// PromptInput represents a natural language prompt for task management
type PromptInput struct {
	Text string `json:"text"`
}

// PromptResponse represents the response to a natural language prompt
type PromptResponse struct {
	Response string `json:"response"`
}
