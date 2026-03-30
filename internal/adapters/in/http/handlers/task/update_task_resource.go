package handlers

import (
	"encoding/json"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/request"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions/codes"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	"github.com/MarcusVNJ/GOTODO/internal/core/usecase"
	"net/http"
)

type UpdateTaskResource struct {
	usecase usecase.IUsecase[*models.Task, struct{}]
}

func NewUpdateTaskResource(usecase usecase.IUsecase[*models.Task, struct{}]) *UpdateTaskResource {
	return &UpdateTaskResource{
		usecase: usecase,
	}
}

func (handler *UpdateTaskResource) Handler(w http.ResponseWriter, r *http.Request) error {
	var updatedTask request.UpdateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		return exceptions.NewBusinessException(codes.UnprocessableEntity)
	}

	task := updatedTask.ToModel()

	_, err := handler.usecase.Execute(r.Context(), task)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responsePayload := map[string]string{
		"message": "Task atualizada com sucesso",
	}

	return json.NewEncoder(w).Encode(responsePayload)
}
