package response

import (
	"github.com/MarcusVNJ/GOTODO/internal/core/enums"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
)

type taskResponse struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Status      enums.Status  `json:"status"`
	Priority    int           `json:"priority"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
}

type taskListItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Priority int    `json:"priority"`
}

type GetTaskResponse struct {
	Body *taskResponse
}

type ListTaskResponse struct {
	Body []taskListItem
}

func NewTaskResponse(task *models.Task) *taskResponse {
	if task == nil {
		return nil
	}

	return &taskResponse{
		ID:          task.ID(),
		Title:       task.Title(),
		Description: task.Description(),
		Status:      task.Status(),
		Priority:    task.Priority(),
		CreatedAt:   task.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   task.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func NewListTaskResponse(tasks []*models.Task) []taskListItem {
	items := make([]taskListItem, 0, len(tasks))
	for _, task := range tasks {
		items = append(items, taskListItem{
			ID:       task.ID(),
			Title:    task.Title(),
			Priority: task.Priority(),
		})
	}
	return items
}