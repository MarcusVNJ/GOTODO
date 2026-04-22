package response

import (
	"github.com/MarcusVNJ/GOTODO/internal/core/enums"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
)

type taskResponse struct {
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Status      enums.Status `json:"status"`
	Priority    int          `json:"priority"`
}

type GetTaskResponse struct {
	Body *taskResponse
}

func NewTaskResponse(task *models.Task) *taskResponse {
	if task == nil {
		return nil
	}

	return &taskResponse{
		Title:       task.Title(),
		Description: task.Description(),
		Status:      task.Status(),
		Priority:    task.Priority(),
	}
}
