package entity

import (
	"github.com/MarcusVNJ/GOTODO/internal/core/enums"
	"time"
)

type TaskEntity struct {
	ID          string     `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Status      enums.Status `db:"status"`
	Priority    int        `db:"priority"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}
