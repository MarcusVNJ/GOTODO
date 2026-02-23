package handlers

import (
	"encoding/json"
    "github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto"
    "github.com/MarcusVNJ/GOTODO/internal/core/base"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions"
	"github.com/MarcusVNJ/GOTODO/internal/core/usecase"
	"net/http"
)

type CreateTaskResource struct {
	usecase base.IUsecase[usecase.REQUEST, usecase.RESPONSE]
}

func NewCreateTaskResource(usecase base.IUsecase[usecase.REQUEST, usecase.RESPONSE]) *CreateTaskResource {
	return &CreateTaskResource{
		usecase: usecase,
	}
}

func (handler *CreateTaskResource) Handler(w http.ResponseWriter, r *http.Request) error {

	var requestDto dto.CreateTaskRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&requestDto); err!= nil {
		return exceptions.NewBusinessException("JSON inválido", err.Error(), exceptions.CodeInvalidData)
	}

	task, err := dto.TaskToModel(requestDto)
	if err!= nil {
		return exceptions.NewBusinessException("Erro inesperado", err.Error(), exceptions.InternalServerError)
	}

	_, err = handler.usecase.Execute(r.Context(), task)
	if err!= nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	responsePayload := map[string]string{
		"message": "Task criada com sucesso",
	}

	return json.NewEncoder(w).Encode(responsePayload)
}
