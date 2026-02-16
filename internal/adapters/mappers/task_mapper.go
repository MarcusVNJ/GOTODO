package mappers

import (
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	"github.com/MarcusVNJ/GOTODO/internal/infrastructure/entity"
	"github.com/rs/xid"
)


func DomainToEntity(task *models.Task) entity.TaskEntity {
	return entity.TaskEntity{
		ID:          task.ID().String(),
		Title:       task.Title(),
		Description: task.Description(),
		Status:      task.Status(),
		Priority:    task.Priority(),
		CreatedAt:   task.CreatedAt(),
		UpdatedAt:   task.UpdatedAt(),
		DeletedAt:   task.DeletedAt(),
	}
}

func EntityToDomain(entity *entity.TaskEntity) (*models.Task, error) {
	id, err := xid.FromString(entity.ID)
	if err!= nil {
		return nil, err
	}

	audit := models.NewAuditInit(id, entity.CreatedAt, entity.UpdatedAt, entity.DeletedAt)

	return models.NewTaskInit(
		audit,
		entity.Title,
		entity.Description,
		entity.Status,
		entity.Priority,
	), nil
}