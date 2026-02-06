package models

import (
	"github.com/rs/xid"
	"time"
)

type Audit struct {
	id        xid.ID
	createdAt time.Time
	updatedAt time.Time
	deletedAt time.Time
}

func NewAudit() Audit {
	timeNow := time.Now().UTC()

	return Audit{
		id:        xid.New(),
		createdAt: timeNow,
		updatedAt: timeNow,
	}
}

func (audit *Audit) UpdatedAudit() {
	audit.updatedAt = time.Now().UTC()
}

func (audit *Audit) GetID() xid.ID {
	return audit.id
}
