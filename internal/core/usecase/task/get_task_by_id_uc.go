package usecase

import (
    "context"
    "github.com/MarcusVNJ/GOTODO/internal/core/base"
    "github.com/MarcusVNJ/GOTODO/internal/core/models"
    "github.com/MarcusVNJ/GOTODO/internal/core/ports"
)

type GetTaskByIdUC struct {
    base.Usecase[string, *models.Task]
    repository repository.TaskRepository
}

func NewGetTaskByIdUC(repository repository.TaskRepository) *GetTaskByIdUC {
    return &GetTaskByIdUC{
        repository: repository,
    }
}

func (usecase *GetTaskByIdUC) Execute(ctx context.Context, id string) (*models.Task, error) {
    return usecase.Call(ctx, id, usecase.getTaskById)
}

func (usecase *GetTaskByIdUC) getTaskById(ctx context.Context, id string) (*models.Task, error) {
    task, err := usecase.repository.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return task, nil
}
