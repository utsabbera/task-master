package api

import "github.com/utsabbera/task-master/core"

func mapTaskToResponse(task *core.Task) Task {
	return Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func mapTasksToResponse(tasks []*core.Task) []Task {
	response := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		response = append(response, mapTaskToResponse(t))
	}
	return response
}
