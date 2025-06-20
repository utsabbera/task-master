package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/utsabbera/task-master/pkg/util"
)

func TestNewTask(t *testing.T) {
	t.Run("should create task with correct properties", func(t *testing.T) {
		title := "Test Task"
		description := "This is a test task"
		priority := PriorityHigh
		tomorrow := time.Now().Add(24 * time.Hour)

		task := NewTask(title, description, &priority, &tomorrow)

		assert.Empty(t, task.ID, "ID should be empty")
		assert.Equal(t, title, task.Title)
		assert.Equal(t, description, task.Description)
		assert.Equal(t, StatusNotStarted, task.Status)
		assert.Equal(t, priority, *task.Priority)
		assert.Equal(t, tomorrow, *task.DueDate)
		assert.Empty(t, task.CreatedAt)
		assert.Empty(t, task.UpdatedAt)
	})

	t.Run("should create task with nil due date", func(t *testing.T) {
		task := NewTask("Test", "Description", util.Ptr(PriorityMedium), nil)

		assert.Nil(t, task.DueDate)
	})

	t.Run("should create task with nil priority", func(t *testing.T) {
		task := NewTask("Test", "Description", nil, nil)

		assert.Nil(t, task.Priority)
	})

	t.Run("should set status to not started", func(t *testing.T) {
		task := NewTask("Test", "Description", nil, nil)

		assert.Equal(t, StatusNotStarted, task.Status)
	})
}

func TestPriority(t *testing.T) {
	t.Run("should have correct priority values", func(t *testing.T) {
		assert.Equal(t, Priority(1), PriorityLow)
		assert.Equal(t, Priority(2), PriorityMedium)
		assert.Equal(t, Priority(3), PriorityHigh)
	})
}

func TestStatus(t *testing.T) {
	t.Run("should have correct status values", func(t *testing.T) {
		assert.Equal(t, Status("NOT_STARTED"), StatusNotStarted)
		assert.Equal(t, Status("IN_PROGRESS"), StatusInProgress)
		assert.Equal(t, Status("COMPLETED"), StatusCompleted)
	})
}
