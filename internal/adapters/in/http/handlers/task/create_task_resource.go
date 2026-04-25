package handlers

import (
	"context"

	"net/http"

	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/request"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/response"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/middlewares"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	"github.com/MarcusVNJ/GOTODO/internal/core/usecase"
	"github.com/danielgtaylor/huma/v2"
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

func (r *CreateTaskResource) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-task",
		Method:        http.MethodPost,
		Path:          "/api/task",
		Summary:       "Criar Tarefa",
		Description:   "Cria uma nova tarefa",
		Tags:          []string{"Tasks"},
		DefaultStatus: http.StatusCreated,
	}, middlewares.HandlerException(r.Handler))
}
