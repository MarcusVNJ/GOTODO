package usecase

import (
	"context"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/mappers"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/out"
	"github.com/MarcusVNJ/GOTODO/internal/core/base"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
)

type REQUEST = models.Task
type RESPONSE = struct{}

type ICreateTaskUC interface {
	base.IUsecase[REQUEST, RESPONSE]
}

type CreateTaskUC struct {
	base.BaseUsecase[REQUEST, RESPONSE]
	repository repository.TaskRepository
}

func NewCreateTaskUC(repository repository.TaskRepository) ICreateTaskUC {
	return &CreateTaskUC{
		repository: repository,
	}
}

func (usecase *CreateTaskUC) Execute(ctx context.Context, request REQUEST) (RESPONSE, error) {
	return usecase.Call(ctx, request, usecase.createTask)
}

func (usecase *CreateTaskUC) createTask(ctx context.Context, request REQUEST) (RESPONSE, error) {

	task := mappers.DomainToEntity(&request)

	err := usecase.repository.Save(ctx, task)
	if err != nil {
		return RESPONSE{}, err
	}

	return struct{}{}, nil
}
