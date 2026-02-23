package repository

import (
	"context"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/out/infrastructure/entity"
	"github.com/MarcusVNJ/GOTODO/internal/core/enums"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	"github.com/rs/xid"
)

type TaskRepository interface {
	Save(context context.Context, task *models.Task) error
	FindByID(context context.Context, id xid.ID) (*models.Task, error)
	FindAll(context context.Context, statusFilter string, minPriority int) ([]*models.Task, error)
	Update(context context.Context, task entity.TaskEntity) error
	Delete(context context.Context, id xid.ID) error
	FindByStatus(context context.Context, status enums.Status) ([]*models.Task, error)
}