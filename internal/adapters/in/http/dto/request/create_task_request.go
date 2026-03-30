package request

import "github.com/MarcusVNJ/GOTODO/internal/core/models"

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
}

func (dto CreateTaskRequest) ToModel() *models.Task {
	return models.NewTask(dto.Title, dto.Description, dto.Priority)
}
