package repository

import (
	"context"
	"github.com/MarcusVNJ/GOTODO/internal/core/enums"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
)

type TaskRepository interface {
	Save(context context.Context, task *models.Task) error
	FindByID(context context.Context, id string) (*models.Task, error)
	ExistByID(context context.Context, id string) (bool, error)
	FindAll(context context.Context, statusFilter string, minPriority int) ([]*models.Task, error)
	Update(context context.Context, task *models.Task) error
	Delete(context context.Context, id string) error
	FindByStatus(context context.Context, status enums.Status) ([]*models.Task, error)
}
