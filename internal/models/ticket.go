package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TicketStatus описывает допустимые статусы заявки.
type TicketStatus string

const (
	// TicketStatusPlanned означает, что заявка запланирована к выполнению.
	TicketStatusPlanned TicketStatus = "PLANNED"
	// TicketStatusInProgress означает, что заявка находится в работе.
	TicketStatusInProgress TicketStatus = "IN_PROGRESS"
	// TicketStatusCompleted означает, что работы по заявке завершены.
	TicketStatusCompleted TicketStatus = "COMPLETED"
	// TicketStatusClosed означает, что заявка закрыта ответственными.
	TicketStatusClosed TicketStatus = "CLOSED"
	// TicketStatusCancelled означает, что заявка отменена.
	TicketStatusCancelled TicketStatus = "CANCELLED"
)

// Ticket описывает заявку на выполнение работ по уборке.
type Ticket struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey"`
	CleaningAreaID uuid.UUID    `gorm:"type:uuid;not null;index"`
	ContractorID   uuid.UUID    `gorm:"type:uuid;not null;index"`
	CreatedByOrgID uuid.UUID    `gorm:"type:uuid;not null;index"`
	Status         TicketStatus `gorm:"type:varchar(20);not null;default:PLANNED" json:"status"`
	PlannedStartAt time.Time
	PlannedEndAt   time.Time
	FactStartAt    *time.Time
	FactEndAt      *time.Time
	Description    string  `gorm:"type:text"`
	PhotoURL       *string `gorm:"type:varchar(255)"`
	Latitude       *float64
	Longitude      *float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TableName возвращает имя таблицы для Ticket.
func (Ticket) TableName() string {
	return "tickets"
}

// BeforeCreate выставляет UUID перед созданием записи.
func (t *Ticket) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
