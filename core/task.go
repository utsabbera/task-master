package core

import (
	"time"
)

// Status represents the current state of a task
type Status string

const (
	// StatusNotStarted indicates a task that hasn't been started yet
	StatusNotStarted Status = "NOT_STARTED"
	// StatusInProgress indicates a task that is currently being worked on
	StatusInProgress Status = "IN_PROGRESS"
	// StatusCompleted indicates a task that has been finished
	StatusCompleted Status = "COMPLETED"
)

// Priority defines the importance level of a task
type Priority int

const (
	// PriorityLow represents the lowest importance level
	PriorityLow Priority = 1
	// PriorityMedium represents the standard importance level
	PriorityMedium Priority = 2
	// PriorityHigh represents the highest importance level
	PriorityHigh Priority = 3
)

// Task represents a single task in the task management system
type Task struct {
	// ID is the unique identifier for the task
	ID string
	// Title is the short name of the task
	Title string
	// Description provides additional details about the task
	Description string
	// Status indicates the current state of the task
	Status Status
	// Priority indicates the importance level of the task
	Priority *Priority
	// CreatedAt stores when the task was created
	CreatedAt time.Time
	// UpdatedAt stores when the task was last modified
	UpdatedAt time.Time
	// DueDate is the optional deadline for the task
	DueDate *time.Time
}

// NewTask creates a new task with the specified properties
// The ID field and timestamps will be managed by the repository
func NewTask(title, description string, priority *Priority, dueDate *time.Time) *Task {
	return &Task{
		Title:       title,
		Description: description,
		Status:      StatusNotStarted,
		Priority:    priority,
		DueDate:     dueDate,
	}
}
