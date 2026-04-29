package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions/codes"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	usecase "github.com/MarcusVNJ/GOTODO/internal/core/usecase/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateTaskUC_Execute_Success(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	uc := usecase.NewUpdateTaskUC(mockRepo)

	task := models.NewTask("Test Task", "Test Description", 1)
	task.SetID("123")

	mockRepo.On("ExistByID", mock.Anything, "123").Return(true, nil)
	mockRepo.On("Update", mock.Anything, task).Return(nil)

	_, err := uc.Execute(context.Background(), task)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTaskUC_Execute_NotFound(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	uc := usecase.NewUpdateTaskUC(mockRepo)

	task := models.NewTask("Test Task", "Test Description", 1)
	task.SetID("123")

	mockRepo.On("ExistByID", mock.Anything, "123").Return(false, nil)

	_, err := uc.Execute(context.Background(), task)

	assert.Error(t, err)
	
	var bizErr *exceptions.BusinessException
	assert.ErrorAs(t, err, &bizErr)
	assert.Equal(t, codes.TaskNotFound.Code(), bizErr.Code)
	
	mockRepo.AssertExpectations(t)
}

func TestUpdateTaskUC_Execute_ExistError(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	uc := usecase.NewUpdateTaskUC(mockRepo)

	task := models.NewTask("Test Task", "Test Description", 1)
	task.SetID("123")

	expectedErr := errors.New("db error")
	mockRepo.On("ExistByID", mock.Anything, "123").Return(false, expectedErr)

	_, err := uc.Execute(context.Background(), task)

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTaskUC_Execute_UpdateError(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	uc := usecase.NewUpdateTaskUC(mockRepo)

	task := models.NewTask("Test Task", "Test Description", 1)
	task.SetID("123")

	expectedErr := errors.New("update error")
	mockRepo.On("ExistByID", mock.Anything, "123").Return(true, nil)
	mockRepo.On("Update", mock.Anything, task).Return(expectedErr)

	_, err := uc.Execute(context.Background(), task)

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	mockRepo.AssertExpectations(t)
}
