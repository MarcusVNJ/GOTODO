package usecase

import (
    "context"
    "github.com/MarcusVNJ/GOTODO/internal/core/base"
    "github.com/MarcusVNJ/GOTODO/internal/core/models"
    "github.com/MarcusVNJ/GOTODO/internal/core/ports"
)

type CreateTaskUC struct {
    base.BaseUsecase[*models.Task, struct{}]
    repository repository.TaskRepository
}

func NewCreateTaskUC(repository repository.TaskRepository) *CreateTaskUC {
    return &CreateTaskUC{
        repository: repository,
    }
}

func (usecase *CreateTaskUC) Execute(ctx context.Context, request *models.Task) (struct{}, error) {
    return usecase.Call(ctx, request, usecase.createTask)
}

func (usecase *CreateTaskUC) createTask(ctx context.Context, request *models.Task) (struct{}, error) {

    err := usecase.repository.Save(ctx, request)
    if err != nil {
        return struct{}{}, err
    }

    return struct{}{}, nil
}
