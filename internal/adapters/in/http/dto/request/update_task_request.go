package request

import "github.com/MarcusVNJ/GOTODO/internal/core/enums"

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
