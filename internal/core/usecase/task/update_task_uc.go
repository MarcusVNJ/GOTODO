package usecase

import (
	"context"
	"github.com/MarcusVNJ/GOTODO/internal/core/base"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions/codes"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	"github.com/MarcusVNJ/GOTODO/internal/core/ports"
)

type UpdateTaskUc struct {
	base.Usecase[*models.Task, struct{}]
	repository repository.TaskRepository
}

func NewUpdateTaskUC(repository repository.TaskRepository) *UpdateTaskUc {
	return &UpdateTaskUc{
		repository: repository,
	}
}

func (usecase *UpdateTaskUc) Execute(ctx context.Context, request *models.Task) (struct{}, error) {
	return usecase.Call(ctx, request, usecase.UpdateTask)
}

func (usecase *UpdateTaskUc) UpdateTask(ctx context.Context, request *models.Task) (struct{}, error) {
	exist, err := usecase.repository.ExistByID(ctx, request.ID())
	if err != nil {
		return struct{}{}, err
	} else if !exist {
		return struct{}{}, exceptions.NewBusinessException(codes.TaskNotFound)
	}

	return struct{}{}, usecase.repository.Update(ctx, request)
}
