package handlers

import (
	"context"
	"net/http"

	taskdto "github.com/MarcusVNJ/GOTODO/internal/app/task/dto"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/request"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/response"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/middlewares"
	"github.com/MarcusVNJ/GOTODO/internal/core/usecase"
	"github.com/danielgtaylor/huma/v2"
)

type UpdateTaskResource struct {
	usecase usecase.IUsecase[taskdto.UpdateTaskCommand, struct{}]
}

func NewUpdateTaskResource(uc usecase.IUsecase[taskdto.UpdateTaskCommand, struct{}]) *UpdateTaskResource {
	return &UpdateTaskResource{usecase: uc}
}

func (r *UpdateTaskResource) Handler(ctx context.Context, input *request.UpdateTaskRequest) (*response.OperationTaskResponse, error) {
	cmd := taskdto.UpdateTaskCommand{
		ID:          input.Body.Id,
		Title:       input.Body.Title,
		Description: input.Body.Description,
		Status:      input.Body.Status,
		Priority:    input.Body.Priority,
	}

	_, err := r.usecase.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &response.OperationTaskResponse{
		Body: response.MessagePayload{
			Message: "Task atualizada com sucesso",
		},
	}, nil
}

func (r *UpdateTaskResource) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "update-task",
		Method:      http.MethodPut,
		Path:        "/api/task",
		Summary:     "Atualizar Tarefa",
		Description: "Atualiza os atributos de uma tarefa",
		Tags:        []string{"Tasks"},
	}, middlewares.HandlerException(r.Handler))
}
