package handlers

import (
	"context"

	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/request"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/response"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	"github.com/MarcusVNJ/GOTODO/internal/core/usecase"
)

type UpdateTaskResource struct {
	usecase usecase.IUsecase[*models.Task, struct{}]
}

func NewUpdateTaskResource(usecase usecase.IUsecase[*models.Task, struct{}]) *UpdateTaskResource {
	return &UpdateTaskResource{
		usecase: usecase,
	}
}

func (r *UpdateTaskResource) Handler(ctx context.Context, input *request.UpdateTaskRequest) (*response.OperationTaskResponse, error) {
	task := input.ToModel()

	_, err := r.usecase.Execute(ctx, task)
	if err != nil {
		return nil, err
	}

	return &response.OperationTaskResponse{
		Body: response.MessagePayload{
			Message: "Task atualizada com sucesso",
		},
	}, nil
}
