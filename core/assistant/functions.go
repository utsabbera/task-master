package assistant

import (
	"context"
	"time"

	"github.com/utsabbera/task-master/core/task"
	"github.com/utsabbera/task-master/pkg/assistant"
	"github.com/utsabbera/task-master/pkg/util"
)

func NewCreateFunction(service task.Service) assistant.Function {
	return assistant.NewFunction(
		"create_task",
		"Create a new task with Title, Description, Priority, DueDate, and Status",
		func(ctx context.Context, params struct {
			Title       string         `json:"title" jsonschema:"description=The title of the task,example=Buy groceries"`
			Description string         `json:"description" jsonschema:"description=The description of the task,example=Milk, eggs, bread"`
			Status      task.Status    `json:"status" jsonschema:"description=The status of the task,enum=NOT_STARTED,enum=IN_PROGRESS,enum=COMPLETED"`
			Priority    *task.Priority `json:"priority,omitempty" jsonschema:"description=The priority of the task,enum=LOW,enum=MEDIUM,enum=HIGH"`
			DueDate     *time.Time     `json:"dueDate,omitempty" jsonschema:"description=The due date of the task,example=2006-01-02T15:04:05Z07:00"`
			// FIXME: Define json schema to send null instead of "" for time
		}) (*task.Task, error) {
			task := &task.Task{
				Title:       params.Title,
				Description: params.Description,
				Status:      params.Status,
				Priority:    params.Priority,
				DueDate:     params.DueDate,
			}

			err := service.Create(task)
			return task, err
		},
	)
}

func NewGetFunction(service task.Service) assistant.Function {
	return assistant.NewFunction(
		"get_task",
		"Get a task by ID",
		func(ctx context.Context, params struct {
			ID string `json:"id" jsonschema:"title=ID,description=The ID of the task,example=123"`
		}) (*task.Task, error) {
			t, err := service.Get(params.ID)
			return t, err
		},
	)
}

func NewListFunction(service task.Service) assistant.Function {
	return assistant.NewFunction(
		"list_tasks",
		"List all tasks",
		func(ctx context.Context, _ struct{}) ([]*task.Task, error) {
			ts, err := service.List()
			return ts, err
		},
	)
}

func NewUpdateFunction(service task.Service) assistant.Function {
	return assistant.NewFunction(
		"update_task",
		"Update a task by ID with Title, Description, Priority, DueDate, and Status",
		func(ctx context.Context, params struct {
			// TODO: Set default values for the optional fields

			ID          string         `json:"id" jsonschema:"description=The id of the task,example=123"`
			Title       string         `json:"title,omitempty" jsonschema:"description=The title of the task,example=Buy groceries"`
			Description string         `json:"description,omitempty" jsonschema:"description=The description of the task,example=Milk, eggs, bread"`
			Status      task.Status    `json:"status,omitempty" jsonschema:"description=The status of the task,enum=NOT_STARTED,enum=IN_PROGRESS,enum=COMPLETED"`
			Priority    *task.Priority `json:"priority,omitempty" jsonschema:"description=The priority of the task,enum=LOW,enum=MEDIUM,enum=HIGH"`
			DueDate     *time.Time     `json:"dueDate,omitempty" jsonschema:"description=The due date of the task,example=2025-06-23T15:04:05Z07:00"`
		}) (*task.Task, error) {
			task := &task.Task{
				Title:       params.Title,
				Description: params.Description,
				Status:      params.Status,
				Priority:    params.Priority,
				DueDate:     params.DueDate,
			}

			t, err := service.Update(params.ID, task)
			return t, err
		},
	)
}

func NewDeleteFunction(service task.Service) assistant.Function {
	return assistant.NewFunction(
		"delete_task",
		"Delete a task by ID",
		func(ctx context.Context, params struct {
			ID string `json:"id" jsonschema:"description=The id of the task,example=123"`
		}) (bool, error) {
			err := service.Delete(params.ID)
			return err == nil, err
		},
	)
}

func NewDateTimeFunction(clock util.Clock) assistant.Function {
	return assistant.NewFunction(
		"get_current_date",
		"Get the current date and time",
		func(ctx context.Context, _ struct{}) (time.Time, error) {
			return clock.Now(), nil
		},
	)
}
