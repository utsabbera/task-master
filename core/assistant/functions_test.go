package assistant_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/utsabbera/task-master/core/assistant"
	coretask "github.com/utsabbera/task-master/core/task"
	pkgassistant "github.com/utsabbera/task-master/pkg/assistant"
	"go.uber.org/mock/gomock"
)

func callFunction[P any](fn pkgassistant.Function, ctx context.Context, params P) (any, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	resp := fn.Call(ctx, string(b))
	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	return resp.Data, nil
}

func TestNewCreateFunction(t *testing.T) {
	t.Run("should generate a valid function definition", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := coretask.NewMockService(ctrl)
		fn := assistant.NewListFunction(mockService)

		def := fn.Definition()
		assert.Equal(t, "create_task", def.Name)
		assert.Equal(t, "Creates a new task", def.Description)

		bytes, err := json.Marshal(def.Parameters)
		require.NoError(t, err)

		assert.Equal(t, "", string(bytes))
	})

	t.Run("should create a task with correct title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := coretask.NewMockService(ctrl)
		fn := assistant.NewCreateFunction(mockService)
		params := struct {
			Title       string
			Description string
			Status      coretask.Status
			Priority    *coretask.Priority
			DueDate     *time.Time
		}{Title: "title", Status: coretask.Status("NOT_STARTED")}

		mockService.EXPECT().Create(gomock.AssignableToTypeOf(&coretask.Task{})).DoAndReturn(
			func(tk *coretask.Task) error {
				if tk.Title != "title" {
					t.Errorf("expected title 'title', got %v", tk.Title)
				}
				return nil
			},
		)
		_, err := callFunction(fn, context.Background(), params)
		assert.NoError(t, err)
	})
}

func TestNewGetFunction(t *testing.T) {
	t.Run("should get a task by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := coretask.NewMockService(ctrl)
		fn := assistant.NewGetFunction(mockService)
		params := struct{ ID string }{ID: "ok"}
		mockService.EXPECT().Get("ok").Return(&coretask.Task{ID: "ok"}, nil)
		res, err := callFunction(fn, context.Background(), params)
		assert.NoError(t, err)
		tk, ok := res.(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "ok", tk["ID"])
	})

	t.Run("should return error if task not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := coretask.NewMockService(ctrl)
		fn := assistant.NewGetFunction(mockService)
		params := struct{ ID string }{ID: "bad"}
		mockService.EXPECT().Get("bad").Return(nil, errors.New("not found"))
		_, err := callFunction(fn, context.Background(), params)
		assert.Error(t, err)
	})
}

func TestNewListFunction(t *testing.T) {
	t.Run("should list all tasks", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := coretask.NewMockService(ctrl)
		fn := assistant.NewListFunction(mockService)
		mockService.EXPECT().List().Return([]*coretask.Task{{ID: "1"}, {ID: "2"}}, nil)
		_, err := callFunction(fn, context.Background(), struct{}{})
		assert.NoError(t, err)
	})
}

func TestNewUpdateFunction(t *testing.T) {
	t.Run("should update a task by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := coretask.NewMockService(ctrl)
		fn := assistant.NewUpdateFunction(mockService)
		params := struct {
			ID          string
			Title       string
			Description string
			Status      coretask.Status
			Priority    *coretask.Priority
			DueDate     *time.Time
		}{ID: "ok", Title: "new"}
		mockService.EXPECT().Update("ok", gomock.AssignableToTypeOf(&coretask.Task{})).Return(&coretask.Task{ID: "ok", Title: "new"}, nil)
		res, err := callFunction(fn, context.Background(), params)
		assert.NoError(t, err)
		tk, ok := res.(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "ok", tk["ID"])
		assert.Equal(t, "new", tk["Title"])
	})

	t.Run("should return error if update fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := coretask.NewMockService(ctrl)
		fn := assistant.NewUpdateFunction(mockService)
		params := struct {
			ID          string
			Title       string
			Description string
			Status      coretask.Status
			Priority    *coretask.Priority
			DueDate     *time.Time
		}{ID: "bad", Title: "new"}
		mockService.EXPECT().Update("bad", gomock.AssignableToTypeOf(&coretask.Task{})).Return(nil, errors.New("not found"))
		_, err := callFunction(fn, context.Background(), params)
		assert.Error(t, err)
	})
}

func TestNewDeleteFunction(t *testing.T) {
	t.Run("should delete a task by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := coretask.NewMockService(ctrl)
		fn := assistant.NewDeleteFunction(mockService)
		params := struct{ ID string }{ID: "ok"}
		mockService.EXPECT().Delete("ok").Return(nil)
		okVal, err := callFunction(fn, context.Background(), params)
		assert.NoError(t, err)
		ok, okType := okVal.(bool)
		assert.True(t, okType)
		assert.True(t, ok)
	})

	t.Run("should return error if delete fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := coretask.NewMockService(ctrl)
		fn := assistant.NewDeleteFunction(mockService)
		params := struct{ ID string }{ID: "bad"}
		mockService.EXPECT().Delete("bad").Return(errors.New("not found"))
		okVal, err := callFunction(fn, context.Background(), params)
		assert.Error(t, err)
		ok, okType := okVal.(bool)
		assert.True(t, okType)
		assert.False(t, ok)
	})
}
