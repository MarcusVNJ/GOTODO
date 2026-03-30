package response

import (
	"github.com/MarcusVNJ/GOTODO/internal/core/enums"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
)

type TaskResponse struct {
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Status      enums.Status `json:"status"`
	Priority    int          `json:"priority"`
}

func NewTaskResponse(task *models.Task) *TaskResponse {
	if task == nil {
		return nil
	}

	return &TaskResponse{
		Title:       task.Title(),
		Description: task.Description(),
		Status:      task.Status(),
		Priority:    task.Priority(),
	}
}
