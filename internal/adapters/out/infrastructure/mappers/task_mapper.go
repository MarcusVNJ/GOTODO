package mappers

import (
	"github.com/MarcusVNJ/GOTODO/internal/adapters/out/infrastructure/entity"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
)

func DomainToEntity(task *models.Task) entity.TaskEntity {
	return entity.TaskEntity{
		ID:          task.ID(),
		Title:       task.Title(),
		Description: task.Description(),
		Status:      task.Status(),
		Priority:    task.Priority(),
		CreatedAt:   task.CreatedAt(),
		UpdatedAt:   task.UpdatedAt(),
		DeletedAt:   task.DeletedAt(),
	}
}

func EntityToDomain(entity *entity.TaskEntity) *models.Task {

	audit := models.NewAuditInit(entity.ID, entity.CreatedAt, entity.UpdatedAt, entity.DeletedAt)

	return models.NewTaskInit(
		audit,
		entity.Title,
		entity.Description,
		entity.Status,
		entity.Priority,
	)
}
