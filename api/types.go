package api

// Task represents a task in the task management system.
type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type TaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	DueDate     string `json:"dueDate,omitempty"`
}
