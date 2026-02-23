package usecase

import (
	"context"
	"github.com/MarcusVNJ/GOTODO/internal/core/base"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	"github.com/MarcusVNJ/GOTODO/internal/core/ports"
)

type REQUEST = *models.Task
type RESPONSE = struct{}

type CreateTaskUC struct {
	base.BaseUsecase[REQUEST, RESPONSE]
	repository repository.TaskRepository
}

func NewCreateTaskUC(repository repository.TaskRepository) *CreateTaskUC {
	return &CreateTaskUC{
		repository: repository,
	}
}

func (usecase *CreateTaskUC) Execute(ctx context.Context, request REQUEST) (RESPONSE, error) {
	return usecase.Call(ctx, request, usecase.createTask)
}

func (usecase *CreateTaskUC) createTask(ctx context.Context, request REQUEST) (RESPONSE, error) {

	err := usecase.repository.Save(ctx, request)
	if err != nil {
		return RESPONSE{}, err
	}

	return RESPONSE{}, nil
}
