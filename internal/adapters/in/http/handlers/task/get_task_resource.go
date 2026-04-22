package handlers

import (
	"context"

	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/request"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/response"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	"github.com/MarcusVNJ/GOTODO/internal/core/usecase"
)

type GetTaskByIdResource struct {
	usecase usecase.IUsecase[string, *models.Task]
}

func NewGetTaskByIdResource(usecase usecase.IUsecase[string, *models.Task]) *GetTaskByIdResource {
	return &GetTaskByIdResource{
		usecase: usecase,
	}
}

func (r *GetTaskByIdResource) Handler(ctx context.Context, input *request.GetTaskRequest) (*response.GetTaskResponse, error) {
	task, err := r.usecase.Execute(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &response.GetTaskResponse{
		Body: response.NewTaskResponse(task),
	}, nil
}
