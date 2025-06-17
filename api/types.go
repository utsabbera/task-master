package api

import (
	"time"

	"github.com/utsabbera/task-master/core"
)

// Task represents a task in the task management system.
type Task struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      core.Status   `json:"status"`
	Priority    *core.Priority `json:"priority"`
	DueDate     *time.Time    `json:"dueDate"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

type TaskInput struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Priority    *core.Priority `json:"priority"`
	DueDate     *time.Time     `json:"dueDate"`
}
