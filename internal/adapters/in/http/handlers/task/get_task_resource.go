package handlers

import (
	"context"
	"net/http"

	taskdto "github.com/MarcusVNJ/GOTODO/internal/app/task/dto"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/request"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/response"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/middlewares"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	"github.com/MarcusVNJ/GOTODO/internal/core/usecase"
	"github.com/danielgtaylor/huma/v2"
)

type GetTaskByIdResource struct {
	usecase usecase.IUsecase[taskdto.GetTaskQuery, *models.Task]
}

func NewGetTaskByIdResource(uc usecase.IUsecase[taskdto.GetTaskQuery, *models.Task]) *GetTaskByIdResource {
	return &GetTaskByIdResource{usecase: uc}
}

func (r *GetTaskByIdResource) Handler(ctx context.Context, input *request.GetTaskRequest) (*response.GetTaskResponse, error) {
	query := taskdto.GetTaskQuery{ID: input.ID}

	task, err := r.usecase.Execute(ctx, query)
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
