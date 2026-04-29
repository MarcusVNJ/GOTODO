package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions/codes"
	usecase "github.com/MarcusVNJ/GOTODO/internal/core/usecase/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteTaskUC_Execute_Success(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	uc := usecase.NewDeleteTaskUC(mockRepo)

	mockRepo.On("ExistByID", mock.Anything, "123").Return(true, nil)
	mockRepo.On("Delete", mock.Anything, "123").Return(nil)

	_, err := uc.Execute(context.Background(), "123")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteTaskUC_Execute_NotFound(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	uc := usecase.NewDeleteTaskUC(mockRepo)

	mockRepo.On("ExistByID", mock.Anything, "123").Return(false, nil)

	_, err := uc.Execute(context.Background(), "123")

	assert.Error(t, err)

	var bizErr *exceptions.BusinessException
	assert.ErrorAs(t, err, &bizErr)
	assert.Equal(t, codes.TaskNotFound.Code(), bizErr.Code)

	mockRepo.AssertExpectations(t)
}

func TestDeleteTaskUC_Execute_ExistError(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	uc := usecase.NewDeleteTaskUC(mockRepo)

	expectedErr := errors.New("db error")
	mockRepo.On("ExistByID", mock.Anything, "123").Return(false, expectedErr)

	_, err := uc.Execute(context.Background(), "123")

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	mockRepo.AssertExpectations(t)
}

func TestDeleteTaskUC_Execute_DeleteError(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	uc := usecase.NewDeleteTaskUC(mockRepo)

	expectedErr := errors.New("delete error")
	mockRepo.On("ExistByID", mock.Anything, "123").Return(true, nil)
	mockRepo.On("Delete", mock.Anything, "123").Return(expectedErr)

	_, err := uc.Execute(context.Background(), "123")

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	mockRepo.AssertExpectations(t)
}
