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

type ListTaskResource struct {
	usecase usecase.IUsecase[taskdto.ListTasksQuery, []*models.Task]
}

func NewListTaskResource(uc usecase.IUsecase[taskdto.ListTasksQuery, []*models.Task]) *ListTaskResource {
	return &ListTaskResource{usecase: uc}
}

func (r *ListTaskResource) Handler(ctx context.Context, input *request.ListTaskRequest) (*response.ListTaskResponse, error) {
	query := taskdto.ListTasksQuery{
		Status:      input.Status,
		MinPriority: input.MinPriority,
	}

	tasks, err := r.usecase.Execute(ctx, query)
	if err != nil {
		return nil, err
	}

	return &response.ListTaskResponse{Body: response.NewListTaskResponse(tasks)}, nil
}

func (r *ListTaskResource) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "list-tasks",
		Method:      http.MethodGet,
		Path:        "/api/task",
		Summary:     "Listar Tarefas",
		Description: "Lista todas as tarefas ativas com informações resumidas",
		Tags:        []string{"Tasks"},
	}, middlewares.HandlerException(r.Handler))
}
