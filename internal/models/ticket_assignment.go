package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TicketAssignment описывает назначение водителя и транспорта на заявку.
type TicketAssignment struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	TicketID   uuid.UUID `gorm:"type:uuid;not null;index"`
	DriverID   uuid.UUID `gorm:"type:uuid;not null;index"`
	VehicleID  uuid.UUID `gorm:"type:uuid;not null;index"`
	AssignedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// TableName возвращает имя таблицы для TicketAssignment.
func (TicketAssignment) TableName() string {
	return "ticket_assignments"
}

// BeforeCreate выставляет UUID и отметку времени назначения.
func (ta *TicketAssignment) BeforeCreate(tx *gorm.DB) (err error) {
	if ta.ID == uuid.Nil {
		ta.ID = uuid.New()
	}
	if ta.AssignedAt.IsZero() {
		ta.AssignedAt = time.Now()
	}
	return nil
}
