package handlers

import (
	"context"

	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/request"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/response"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	"github.com/MarcusVNJ/GOTODO/internal/core/usecase"
)

type CreateTaskResource struct {
	usecase usecase.IUsecase[*models.Task, struct{}]
}

func NewCreateTaskResource(usecase usecase.IUsecase[*models.Task, struct{}]) *CreateTaskResource {
	return &CreateTaskResource{
		usecase: usecase,
	}
}

func (r *CreateTaskResource) Handler(ctx context.Context, input *request.CreateTaskRequest) (*response.OperationTaskResponse, error) {
	task := input.ToModel()

	_, err := r.usecase.Execute(ctx, task)
	if err != nil {
		return nil, err
	}

	return &response.OperationTaskResponse{
		Body: response.MessagePayload{
			Message: "Task criada com sucesso",
			Id:      task.ID(),
		},
	}, nil
}
