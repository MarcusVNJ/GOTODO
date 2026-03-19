package handlers

import (
    "encoding/json"
    "github.com/MarcusVNJ/GOTODO/internal/adapters/in/http/dto"
    "github.com/MarcusVNJ/GOTODO/internal/core/exceptions"
    "github.com/MarcusVNJ/GOTODO/internal/core/models"
    "github.com/MarcusVNJ/GOTODO/internal/core/usecase"
    "net/http"
)

type CreateTaskResource struct {
    usecase usecase.IUsecase[*models.Task, struct{}]
}

func NewCreateTaskResource(usecase usecase.IUsecase[*models.Task, struct{}]) *CreateTaskResource {
    return &CreateTaskResource{
        usecase: usecase,
    }
}

func (handler *CreateTaskResource) Handler(w http.ResponseWriter, r *http.Request) error {

    var requestDto dto.CreateTaskRequestDTO
    if err := json.NewDecoder(r.Body).Decode(&requestDto); err != nil {
        return exceptions.NewBusinessException("JSON inválido", err.Error(), exceptions.UnprocessableEntity)
    }

    task, err := dto.TaskToModel(requestDto)
    if err != nil {
        return err
    }

    _, err = handler.usecase.Execute(r.Context(), task)
    if err != nil {
        return err
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)

    responsePayload := map[string]string{
        "message": "Task criada com sucesso",
    }

    return json.NewEncoder(w).Encode(responsePayload)
}
