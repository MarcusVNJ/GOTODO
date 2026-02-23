package dto

import "github.com/MarcusVNJ/GOTODO/internal/core/models"

type CreateTaskRequestDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
}

func TaskToModel(dto CreateTaskRequestDTO) (*models.Task, error) {
	return models.NewTask(dto.Title, dto.Description, dto.Priority)
}
