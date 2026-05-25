package handlers

import (
	"context"
	"net/http"

	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/request"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/response"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/middlewares"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions/codes"
	"github.com/MarcusVNJ/GOTODO/internal/core/usecase"
	"github.com/danielgtaylor/huma/v2"
)

type DeleteTaskResource struct {
	usecase usecase.IUsecase[string, struct{}]
}

func NewDeleteTaskResource(uc usecase.IUsecase[string, struct{}]) *DeleteTaskResource {
	return &DeleteTaskResource{usecase: uc}
}

func (r *DeleteTaskResource) Handler(ctx context.Context, input *request.DeleteTaskRequest) (*response.OperationTaskResponse, error) {
	if err := validateId(input.ID); err != nil {
		return nil, err
	}

	_, err := r.usecase.Execute(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &response.OperationTaskResponse{
		Body: response.MessagePayload{
			Message: "Task deletada com sucesso",
		},
	}, nil
}

func (r *DeleteTaskResource) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "delete-task",
		Method:        http.MethodDelete,
		Path:          "/api/task/{id}",
		Summary:       "Excluir Tarefa",
		Description:   "Remove uma tarefa por ID",
		Tags:          []string{"Tasks"},
		DefaultStatus: http.StatusNoContent,
	}, middlewares.HandlerException(r.Handler))
}

func validateId(id string) error {
	if id == "" {
		return exceptions.NewBusinessException(codes.IdInvalid)
	}
	return nil
}
