package models

import (
	"github.com/rs/xid"
	"time"
)

type Audit struct {
	id        string
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func NewAudit() Audit {
	timeNow := time.Now().UTC()

	return Audit{
		id:        xid.New().String(),
		createdAt: timeNow,
		updatedAt: timeNow,
	}
}

func NewAuditInit(id string, created, updated time.Time, deleted *time.Time) Audit {
	return Audit{
		id:        id,
		createdAt: created,
		updatedAt: updated,
		deletedAt: deleted,
	}
}

func (audit *Audit) UpdatedAudit() {
	audit.updatedAt = time.Now().UTC()
}

func (audit Audit) SetID(id string)       { audit.id = id }
func (audit Audit) ID() string            { return audit.id }
func (audit Audit) CreatedAt() time.Time  { return audit.createdAt }
func (audit Audit) UpdatedAt() time.Time  { return audit.updatedAt }
func (audit Audit) DeletedAt() *time.Time { return audit.deletedAt }
