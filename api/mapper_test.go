package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/utsabbera/task-master/core"
)

func TestMapTaskToResponse(t *testing.T) {
	due := time.Now().Add(24 * time.Hour)
	priority := core.PriorityHigh
	coreTask := &core.Task{
		ID:          "1",
		Title:       "Test Task",
		Description: "Test Desc",
		Status:      core.StatusInProgress,
		Priority:    &priority,
		DueDate:     &due,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	resp := mapTaskToResponse(coreTask)

	assert.Equal(t, coreTask.ID, resp.ID)
	assert.Equal(t, coreTask.Title, resp.Title)
	assert.Equal(t, coreTask.Description, resp.Description)
	assert.Equal(t, coreTask.Status, resp.Status)
	assert.Equal(t, coreTask.Priority, resp.Priority)
	assert.Equal(t, coreTask.DueDate, resp.DueDate)
	assert.Equal(t, coreTask.CreatedAt, resp.CreatedAt)
	assert.Equal(t, coreTask.UpdatedAt, resp.UpdatedAt)
}

func TestMapTasksToResponse(t *testing.T) {
	coreTasks := []*core.Task{
		{ID: "1", Title: "A"},
		{ID: "2", Title: "B"},
	}
	resp := mapTasksToResponse(coreTasks)
	assert.Equal(t, len(coreTasks), len(resp))
	for i := range coreTasks {
		assert.Equal(t, coreTasks[i].ID, resp[i].ID)
		assert.Equal(t, coreTasks[i].Title, resp[i].Title)
	}
}
