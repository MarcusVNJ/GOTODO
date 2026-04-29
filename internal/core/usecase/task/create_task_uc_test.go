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

func TestCreateTaskUC_Execute_Success(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	uc := usecase.NewCreateTaskUC(mockRepo)

	task := models.NewTask("Test Task", "Test Description", 1)
	
	mockRepo.On("Save", mock.Anything, task).Return(nil)

	_, err := uc.Execute(context.Background(), task)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateTaskUC_Execute_RepositoryError(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	uc := usecase.NewCreateTaskUC(mockRepo)

	task := models.NewTask("Test Task", "Test Description", 1)
	expectedErr := errors.New("repository error")

	mockRepo.On("Save", mock.Anything, task).Return(expectedErr)

	_, err := uc.Execute(context.Background(), task)

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	mockRepo.AssertExpectations(t)
}
