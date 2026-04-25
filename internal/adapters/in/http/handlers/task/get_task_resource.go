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

func (r *GetTaskByIdResource) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-task-by-id",
		Method:      http.MethodGet,
		Path:        "/api/task/{id}",
		Summary:     "Buscar Tarefa",
		Description: "Busca uma tarefa por ID",
		Tags:        []string{"Tasks"},
	}, middlewares.HandlerException(r.Handler))
}
