package handlers

import (
	"context"

	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/request"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/response"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions/codes"
	"github.com/MarcusVNJ/GOTODO/internal/core/usecase"
)

type DeleteTaskResource struct {
	usecase usecase.IUsecase[string, struct{}]
}

func NewDeleteTaskResource(usecase usecase.IUsecase[string, struct{}]) *DeleteTaskResource {
	return &DeleteTaskResource{
		usecase: usecase,
	}
}

func (r *DeleteTaskResource) Handler(ctx context.Context, input *request.DeleteTaskRequest) (*response.OperationTaskResponse, error) {
	err := validateId(input.ID)
	if err != nil {
		return nil, err
	}

	_, err = r.usecase.Execute(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &response.OperationTaskResponse{
		Body: response.MessagePayload{
			Message: "Task deletada com sucesso",
		},
	}, nil
}

func validateId(id string) error {
	if id == "" {
		return exceptions.NewBusinessException(codes.IdInvalid)
	}
	return nil
}
