package handlers

import (
	"encoding/json"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto/response"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	"github.com/MarcusVNJ/GOTODO/internal/core/usecase"
	"net/http"
)

type GetTaskByIdResource struct {
	usecase usecase.IUsecase[string, *models.Task]
}

func NewGetTaskByIdResource(usecase usecase.IUsecase[string, *models.Task]) *GetTaskByIdResource {
	return &GetTaskByIdResource{
		usecase: usecase,
	}
}

func (handler *GetTaskByIdResource) Handler(w http.ResponseWriter, r *http.Request) error {
	taskId := r.PathValue("id")

	task, err := handler.usecase.Execute(r.Context(), taskId)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	taskResponse := response.NewTaskResponse(task)

	return json.NewEncoder(w).Encode(taskResponse)
}
