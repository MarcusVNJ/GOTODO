package handlers

import (
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions"
	"github.com/MarcusVNJ/GOTODO/internal/core/usecase"
	"net/http"
)

type DeleteTaskResource struct {
	usecase usecase.IUsecase[string, struct{}]
}

func NewDeleteTaskResource(usecase usecase.IUsecase[string, struct{}]) *DeleteTaskResource {
	return &DeleteTaskResource{
		usecase: usecase,
	}
}

func (handler *DeleteTaskResource) Handler(w http.ResponseWriter, r *http.Request) error {
	taskId := r.PathValue("id")

	err := validateId(taskId)
	if err != nil {
		return err
	}

	_, err = handler.usecase.Execute(r.Context(), taskId)
	if err != nil {
		return err
	}

}

func validateId(id string) error {
	if id == "" {
		return exceptions.NewBusinessException("ID inválido", nil, exceptions.BadRequest)
	}
	return nil
}
