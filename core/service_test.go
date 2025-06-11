package task

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(task *Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockRepository) FindByID(id string) (*Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Task), args.Error(1)
}

func (m *MockRepository) FindAll() ([]*Task, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Task), args.Error(1)
}

func (m *MockRepository) Update(task *Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateTask(t *testing.T) {
	t.Run("should create task with valid data", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		
		mockRepo.On("Create", mock.AnythingOfType("*task.Task")).Return(nil).Run(func(args mock.Arguments) {
			task := args.Get(0).(*Task)
			task.ID = "TEST-ID"
		})
		
		task, err := service.CreateTask("Test Task", "Description", PriorityMedium, nil)
		
		assert.NoError(t, err)
		assert.Equal(t, "TEST-ID", task.ID)
		assert.Equal(t, "Test Task", task.Title)
		assert.Equal(t, "Description", task.Description)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("should return error with empty title", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		
		task, err := service.CreateTask("", "Description", PriorityMedium, nil)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty")
		assert.Nil(t, task)
		mockRepo.AssertNotCalled(t, "Create")
	})
	
	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		repoErr := errors.New("database error")
		
		mockRepo.On("Create", mock.AnythingOfType("*task.Task")).Return(repoErr)
		
		task, err := service.CreateTask("Test Task", "Description", PriorityMedium, nil)

		assert.Error(t, err)
		assert.Nil(t, task)
		assert.ErrorContains(t, err, "creating task")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetTask(t *testing.T) {
	t.Run("should get task with valid ID", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		expectedTask := &Task{ID: "TEST-ID", Title: "Test Task"}
		
		mockRepo.On("FindByID", "TEST-ID").Return(expectedTask, nil)
		
		task, err := service.GetTask("TEST-ID")
		
		assert.NoError(t, err)
		assert.Equal(t, expectedTask, task)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("should return error with empty ID", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		
		task, err := service.GetTask("")
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty")
		assert.Nil(t, task)
		mockRepo.AssertNotCalled(t, "FindByID")
	})
	
	t.Run("should return repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		
		mockRepo.On("FindByID", "TEST-ID").Return(nil, ErrTaskNotFound)
		
		task, err := service.GetTask("TEST-ID")
		
		assert.Error(t, err)
		assert.Nil(t, task)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllTasks(t *testing.T) {
	t.Run("should return all tasks from repository", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		tasks := []*Task{
			{ID: "1", Title: "Task 1"},
			{ID: "2", Title: "Task 2"},
		}
		
		mockRepo.On("FindAll").Return(tasks, nil)
		
		result, err := service.GetAllTasks()
		
		assert.NoError(t, err)
		assert.Equal(t, tasks, result)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("should return repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		repoErr := errors.New("database error")
		
		mockRepo.On("FindAll").Return(nil, repoErr)
		
		tasks, err := service.GetAllTasks()
		
		assert.Error(t, err)
		assert.Nil(t, tasks)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateTask(t *testing.T) {
	t.Run("should update valid task", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		task := &Task{ID: "TEST-ID", Title: "Test Task"}
		
		mockRepo.On("Update", task).Return(nil)
		
		err := service.UpdateTask(task)
		
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("should call repository update", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		task := &Task{
			ID: "TEST-ID", 
			Title: "Test Task",
		}
		
		mockRepo.On("Update", task).Return(nil)
		
		err := service.UpdateTask(task)
		
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("should return error for nil task", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		
		err := service.UpdateTask(nil)
		
		assert.ErrorIs(t, err, ErrInvalidTask)
		mockRepo.AssertNotCalled(t, "Update")
	})
	
	t.Run("should return repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		task := &Task{ID: "TEST-ID", Title: "Test Task"}
		repoErr := errors.New("database error")
		
		mockRepo.On("Update", mock.AnythingOfType("*task.Task")).Return(repoErr)
		
		err := service.UpdateTask(task)
		
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteTask(t *testing.T) {
	t.Run("should delete task with valid ID", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		
		mockRepo.On("Delete", "TEST-ID").Return(nil)
		
		err := service.DeleteTask("TEST-ID")
		
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("should return error with empty ID", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		
		err := service.DeleteTask("")
		
		assert.Error(t, err)
		mockRepo.AssertNotCalled(t, "Delete")
	})
	
	t.Run("should return repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		repoErr := errors.New("database error")
		
		mockRepo.On("Delete", "TEST-ID").Return(repoErr)
		
		err := service.DeleteTask("TEST-ID")
		
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateTaskStatus(t *testing.T) {
	t.Run("should update task status", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		task := &Task{ID: "TEST-ID", Status: StatusNotStarted}
		
		mockRepo.On("FindByID", "TEST-ID").Return(task, nil)
		mockRepo.On("Update", mock.AnythingOfType("*task.Task")).Return(nil).Run(func(args mock.Arguments) {
			updatedTask := args.Get(0).(*Task)
			assert.Equal(t, StatusCompleted, updatedTask.Status)
		})
		
		err := service.UpdateTaskStatus("TEST-ID", StatusCompleted)
		
		assert.NoError(t, err)
		assert.Equal(t, StatusCompleted, task.Status)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("should call repository update after status change", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		task := &Task{
			ID: "TEST-ID", 
			Status: StatusNotStarted,
		}
		
		mockRepo.On("FindByID", "TEST-ID").Return(task, nil)
		mockRepo.On("Update", task).Return(nil)
		
		err := service.UpdateTaskStatus("TEST-ID", StatusInProgress)
		
		assert.NoError(t, err)
		assert.Equal(t, StatusInProgress, task.Status)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("should return error when finding task fails", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		repoErr := errors.New("database error")
		
		mockRepo.On("FindByID", "TEST-ID").Return(nil, repoErr)
		
		err := service.UpdateTaskStatus("TEST-ID", StatusCompleted)
		
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Update")
	})
	
	t.Run("should return error when update fails", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)
		task := &Task{ID: "TEST-ID", Status: StatusNotStarted}
		repoErr := errors.New("database error")
		
		mockRepo.On("FindByID", "TEST-ID").Return(task, nil)
		mockRepo.On("Update", mock.AnythingOfType("*task.Task")).Return(repoErr)
		
		err := service.UpdateTaskStatus("TEST-ID", StatusCompleted)
		
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
