package handlers

import (
	"encoding/json"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions/codes"
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	responsePayload := map[string]string{
		"message": "Task criada com sucesso",
	}

	return json.NewEncoder(w).Encode(responsePayload)
}

func validateId(id string) error {
	if id == "" {
		return exceptions.NewBusinessException(codes.IdInvalid)
	}
	return nil
}
