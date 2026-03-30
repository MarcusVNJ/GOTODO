package request

import (
	"github.com/MarcusVNJ/GOTODO/internal/core/enums"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
)

type UpdateTaskRequest struct {
	Id          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      enums.Status `json:"status"`
	Priority    int          `json:"priority"`
}

func (dto UpdateTaskRequest) ToModel() *models.Task {

	task := models.NewTaskWithoutAudit(
		dto.Title,
		dto.Description,
		dto.Priority,
	)

	task.Audit.SetID(dto.Id)
	return task
}
