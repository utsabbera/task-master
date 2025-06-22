package task

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/utsabbera/task-master/pkg/idgen"
	"github.com/utsabbera/task-master/pkg/match"
	"github.com/utsabbera/task-master/pkg/util"
	"go.uber.org/mock/gomock"
)

func TestService_Create(t *testing.T) {
	t.Run("should create task with valid data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clock := util.NewMockClock(ctrl)
		mockRepo := NewMockRepository(ctrl)
		mockIdGen := idgen.NewMockGenerator(ctrl)
		service := NewService(mockRepo, mockIdGen, clock)

		priority := PriorityMedium
		due := time.Now().Add(24 * time.Hour).Truncate(time.Second)

		mockIdGen.EXPECT().Next().Return("TEST-ID")
		createTime := time.Now().Truncate(time.Second)
		clock.EXPECT().Now().Return(createTime)

		task := &Task{
			Title:       "Test Task",
			Description: "Description",
			Status:      StatusInProgress,
			Priority:    &priority,
			DueDate:     &due,
		}

		mockRepo.EXPECT().Create(match.PtrTo(&Task{
			ID:          "TEST-ID",
			Title:       "Test Task",
			Description: "Description",
			Status:      StatusInProgress,
			Priority:    &priority,
			DueDate:     &due,
			CreatedAt:   createTime,
			UpdatedAt:   createTime,
		})).Return(nil)

		err := service.Create(task)

		assert.NoError(t, err)
		assert.Equal(t, "TEST-ID", task.ID)
		assert.Equal(t, "Test Task", task.Title)
		assert.Equal(t, "Description", task.Description)
		assert.Equal(t, &priority, task.Priority)
		assert.WithinDuration(t, due, *task.DueDate, time.Second)
	})

	t.Run("should create task without priority", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clock := util.NewMockClock(ctrl)
		mockRepo := NewMockRepository(ctrl)
		mockIdGen := idgen.NewMockGenerator(ctrl)
		service := NewService(mockRepo, mockIdGen, clock)

		due := time.Now().Add(24 * time.Hour).Truncate(time.Second)

		mockIdGen.EXPECT().Next().Return("TEST-ID")
		createTime := time.Now().Truncate(time.Second)
		clock.EXPECT().Now().Return(createTime)

		task := &Task{
			Title:       "No Priority",
			Description: "Description",
			Status:      StatusInProgress,
			DueDate:     &due,
		}

		mockRepo.EXPECT().Create(match.PtrTo(&Task{
			ID:          "TEST-ID",
			Title:       "No Priority",
			Description: "Description",
			Status:      StatusInProgress,
			Priority:    nil,
			DueDate:     &due,
			CreatedAt:   createTime,
			UpdatedAt:   createTime,
		})).Return(nil)

		err := service.Create(task)

		assert.NoError(t, err)
		assert.Equal(t, "TEST-ID", task.ID)
		assert.Nil(t, task.Priority)
	})

	t.Run("should create task without due date", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clock := util.NewMockClock(ctrl)
		mockRepo := NewMockRepository(ctrl)
		mockIdGen := idgen.NewMockGenerator(ctrl)
		service := NewService(mockRepo, mockIdGen, clock)

		priority := PriorityMedium
		createTime := time.Now().Truncate(time.Second)
		clock.EXPECT().Now().Return(createTime)
		mockIdGen.EXPECT().Next().Return("TEST-ID")

		task := &Task{
			Title:       "No DueDate",
			Description: "Description",
			Status:      StatusInProgress,
			Priority:    &priority,
		}

		mockRepo.EXPECT().Create(match.PtrTo(&Task{
			ID:          "TEST-ID",
			Title:       "No DueDate",
			Description: "Description",
			Status:      StatusInProgress,
			Priority:    &priority,
			DueDate:     nil,
			CreatedAt:   createTime,
			UpdatedAt:   createTime,
		})).Return(nil)

		err := service.Create(task)

		assert.NoError(t, err)
		assert.Equal(t, "TEST-ID", task.ID)
		assert.Nil(t, task.DueDate)
	})

	t.Run("should set status to not started when not passed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clock := util.NewMockClock(ctrl)
		mockRepo := NewMockRepository(ctrl)
		mockIdGen := idgen.NewMockGenerator(ctrl)
		service := NewService(mockRepo, mockIdGen, clock)

		mockIdGen.EXPECT().Next().Return("TEST-ID")
		createTime := time.Now().Truncate(time.Second)
		clock.EXPECT().Now().Return(createTime)

		task := &Task{
			Title:       "Task Without Status",
			Description: "Description",
			Priority:    nil,
			DueDate:     nil,
		}

		mockRepo.EXPECT().Create(match.PtrTo(&Task{
			ID:          "TEST-ID",
			Title:       "Task Without Status",
			Description: "Description",
			Status:      StatusNotStarted,
			Priority:    nil,
			DueDate:     nil,
			CreatedAt:   createTime,
			UpdatedAt:   createTime,
		})).Return(nil)

		err := service.Create(task)

		assert.NoError(t, err)
		assert.Equal(t, StatusNotStarted, task.Status)
	})

	t.Run("should return error when repository create fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clock := util.NewMockClock(ctrl)
		mockRepo := NewMockRepository(ctrl)
		mockIdGen := idgen.NewMockGenerator(ctrl)
		service := NewService(mockRepo, mockIdGen, clock)

		createTime := time.Now().Truncate(time.Second)
		clock.EXPECT().Now().Return(createTime)
		mockIdGen.EXPECT().Next().Return("TEST-ID")

		task := &Task{
			Title:       "Test Task",
			Description: "Description",
		}

		mockRepo.EXPECT().Create(match.PtrTo(&Task{
			ID:          "TEST-ID",
			Title:       "Test Task",
			Description: "Description",
			Status:      StatusNotStarted,
			Priority:    nil,
			DueDate:     nil,
			CreatedAt:   createTime,
			UpdatedAt:   createTime,
		})).Return(errors.New("repository error"))

		err := service.Create(task)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "repository error")
	})
}

func TestService_Get(t *testing.T) {
	t.Run("should get task by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clock := util.NewMockClock(ctrl)
		mockRepo := NewMockRepository(ctrl)
		mockIdGen := idgen.NewMockGenerator(ctrl)
		service := NewService(mockRepo, mockIdGen, clock)

		mockRepo.EXPECT().
			Get("TEST-ID").
			Return(&Task{ID: "TEST-ID", Title: "Test Task"}, nil)

		task, err := service.Get("TEST-ID")

		assert.NoError(t, err)
		assert.Equal(t, "TEST-ID", task.ID)
		assert.Equal(t, "Test Task", task.Title)
	})

	t.Run("should return error when task not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clock := util.NewMockClock(ctrl)
		mockRepo := NewMockRepository(ctrl)
		mockIdGen := idgen.NewMockGenerator(ctrl)
		service := NewService(mockRepo, mockIdGen, clock)

		mockRepo.EXPECT().
			Get("UNKNOWN").
			Return(nil, errors.New("not found"))

		task, err := service.Get("UNKNOWN")

		assert.Error(t, err)
		assert.Nil(t, task)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestService_List(t *testing.T) {
	t.Run("should list all tasks", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clock := util.NewMockClock(ctrl)
		mockRepo := NewMockRepository(ctrl)
		mockIdGen := idgen.NewMockGenerator(ctrl)
		service := NewService(mockRepo, mockIdGen, clock)

		tasks := []*Task{
			{ID: "1", Title: "Task 1"},
			{ID: "2", Title: "Task 2"},
		}

		mockRepo.EXPECT().
			List().
			Return(tasks, nil)

		result, err := service.List()

		assert.NoError(t, err)
		assert.Equal(t, tasks, result)
		assert.Len(t, result, 2)
	})

	t.Run("should return empty list when no tasks", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clock := util.NewMockClock(ctrl)
		mockRepo := NewMockRepository(ctrl)
		mockIdGen := idgen.NewMockGenerator(ctrl)
		service := NewService(mockRepo, mockIdGen, clock)

		mockRepo.EXPECT().
			List().
			Return([]*Task{}, nil)

		result, err := service.List()

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clock := util.NewMockClock(ctrl)
		mockRepo := NewMockRepository(ctrl)
		mockIdGen := idgen.NewMockGenerator(ctrl)
		service := NewService(mockRepo, mockIdGen, clock)

		mockRepo.EXPECT().
			List().
			Return(nil, errors.New("repository error"))

		result, err := service.List()

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "repository error")
	})
}

func TestService_Update(t *testing.T) {
	t.Run("should update task fields", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clock := util.NewMockClock(ctrl)
		mockRepo := NewMockRepository(ctrl)
		mockIdGen := idgen.NewMockGenerator(ctrl)
		service := NewService(mockRepo, mockIdGen, clock)

		id := "TEST-ID"
		existing := &Task{
			ID:          id,
			Title:       "Old Title",
			Description: "Old Desc",
			Status:      StatusNotStarted,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		patch := &Task{
			Title:       "New Title",
			Description: "New Desc",
		}

		updateTime := existing.UpdatedAt.Add(1 * time.Minute)

		updated := &Task{
			ID:          id,
			Title:       "New Title",
			Description: "New Desc",
			Status:      existing.Status,
			CreatedAt:   existing.CreatedAt,
			UpdatedAt:   updateTime,
		}

		mockRepo.EXPECT().Get(id).Return(existing, nil)
		mockRepo.EXPECT().Update(match.PtrTo(updated)).Return(nil)
		clock.EXPECT().Now().Return(updateTime)

		result, err := service.Update(id, patch)

		assert.NoError(t, err)
		assert.Equal(t, "New Title", result.Title)
		assert.Equal(t, "New Desc", result.Description)
	})

	t.Run("should return error if repo.Get fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clock := util.NewMockClock(ctrl)
		mockRepo := NewMockRepository(ctrl)
		mockIdGen := idgen.NewMockGenerator(ctrl)
		service := NewService(mockRepo, mockIdGen, clock)

		mockRepo.EXPECT().Get("BAD-ID").Return(nil, errors.New("not found"))

		patch := &Task{Title: "Patch"}
		result, err := service.Update("BAD-ID", patch)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should return error if repo.Update fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clock := util.NewMockClock(ctrl)
		mockRepo := NewMockRepository(ctrl)
		mockIdGen := idgen.NewMockGenerator(ctrl)
		service := NewService(mockRepo, mockIdGen, clock)

		id := "TEST-ID"
		existing := &Task{ID: id, Title: "Old Title", Status: StatusNotStarted}
		patch := &Task{Title: "New Title"}
		updateTime := existing.UpdatedAt.Add(1 * time.Minute)

		updated := &Task{
			ID:        id,
			Title:     "New Title",
			Status:    existing.Status,
			CreatedAt: existing.CreatedAt,
			UpdatedAt: updateTime,
		}

		mockRepo.EXPECT().Get(id).Return(existing, nil)
		mockRepo.EXPECT().Update(updated).Return(errors.New("update error"))
		clock.EXPECT().Now().Return(updateTime)

		result, err := service.Update(id, patch)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "update error")
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("should delete task", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clock := util.NewMockClock(ctrl)
		mockRepo := NewMockRepository(ctrl)
		mockIdGen := idgen.NewMockGenerator(ctrl)
		service := NewService(mockRepo, mockIdGen, clock)

		mockRepo.EXPECT().
			Delete("TEST-ID").
			Return(nil)

		err := service.Delete("TEST-ID")

		assert.NoError(t, err)
	})

	t.Run("should return error when delete fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		clock := util.NewMockClock(ctrl)
		mockRepo := NewMockRepository(ctrl)
		mockIdGen := idgen.NewMockGenerator(ctrl)
		service := NewService(mockRepo, mockIdGen, clock)

		mockRepo.EXPECT().
			Delete("TEST-ID").
			Return(errors.New("delete error"))

		err := service.Delete("TEST-ID")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete error")
	})
}
