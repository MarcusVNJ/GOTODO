package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	usecase "github.com/MarcusVNJ/GOTODO/internal/core/usecase/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetTaskByIdUC_Execute_Success(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	uc := usecase.NewGetTaskByIdUC(mockRepo)

	task := models.NewTask("Test Task", "Test Description", 1)
	task.SetID("123")

	mockRepo.On("FindByID", mock.Anything, "123").Return(task, nil)

	result, err := uc.Execute(context.Background(), "123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "123", result.ID())
	mockRepo.AssertExpectations(t)
}

func TestGetTaskByIdUC_Execute_NotFound(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	uc := usecase.NewGetTaskByIdUC(mockRepo)

	expectedErr := errors.New("task not found")

	mockRepo.On("FindByID", mock.Anything, "123").Return((*models.Task)(nil), expectedErr)

	result, err := uc.Execute(context.Background(), "123")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)
	mockRepo.AssertExpectations(t)
}
