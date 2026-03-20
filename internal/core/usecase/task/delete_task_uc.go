package usecase

import (
	"context"
	"github.com/MarcusVNJ/GOTODO/internal/core/base"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions/codes"
	"github.com/MarcusVNJ/GOTODO/internal/core/ports"
)

type DeleteTaskUC struct {
	base.BaseUsecase[string, struct{}]
	repository repository.TaskRepository
}

func NewDeleteTaskUC(repository repository.TaskRepository) *DeleteTaskUC {
	return &DeleteTaskUC{
		repository: repository,
	}
}

func (usecase *DeleteTaskUC) Execute(ctx context.Context, id string) (struct{}, error) {
	return usecase.Call(ctx, id, usecase.deleteTask)
}

func (usecase *DeleteTaskUC) deleteTask(ctx context.Context, id string) (struct{}, error) {
	exist, err := usecase.repository.ExistByID(ctx, id)
	if err != nil {
		return struct{}{}, err
	}

	if !exist {
		return struct{}{}, exceptions.NewBusinessException(codes.TaskNotFound)
	}

	err = usecase.repository.Delete(ctx, id)
	if err != nil {
		return struct{}{}, err
	}
	return struct{}{}, nil
}
