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

type CreateTaskResource struct {
	usecase usecase.IUsecase[taskdto.CreateTaskCommand, taskdto.CreateTaskResult]
}

func NewCreateTaskResource(uc usecase.IUsecase[taskdto.CreateTaskCommand, taskdto.CreateTaskResult]) *CreateTaskResource {
	return &CreateTaskResource{usecase: uc}
}

func (r *CreateTaskResource) Handler(ctx context.Context, input *request.CreateTaskRequest) (*response.OperationTaskResponse, error) {
	cmd := taskdto.CreateTaskCommand{
		Title:       input.Body.Title,
		Description: input.Body.Description,
		Priority:    input.Body.Priority,
	}

	result, err := r.usecase.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &response.OperationTaskResponse{
		Body: response.MessagePayload{
			Message: "Task criada com sucesso",
			Id:      result.ID,
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
