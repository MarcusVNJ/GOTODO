package request

import (
	"time"

	"github.com/MarcusVNJ/GOTODO/internal/core/enums"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
)

type updateTaskPayload struct {
	Id          string       `json:"id" minLength:"1" description:"ID da task" example:"1234"`
	Title       string       `json:"title" minLength:"1" maxLength:"150"`
	Description string       `json:"description" maxLength:"500"`
	Status      enums.Status `json:"status"`
	Priority    int          `json:"priority" minimum:"1" maximum:"5"`
}

type UpdateTaskRequest struct {
	Body updateTaskPayload
}

func (dto UpdateTaskRequest) ToModel() *models.Task {

	audit := models.NewAuditInit(dto.Body.Id, time.Time{}, time.Now().UTC(), nil)

	return models.NewTaskInit(
		audit,
		dto.Body.Title,
		dto.Body.Description,
		dto.Body.Status,
		dto.Body.Priority,
	)
}
