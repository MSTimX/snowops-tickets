package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ticket-service/internal/models"
)

// TicketRepository инкапсулирует операции доступа к данным заявок.
type TicketRepository struct {
	db *gorm.DB
}

// NewTicketRepository создаёт новый экземпляр TicketRepository.
func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

// Create записывает новую заявку в базу данных.
func (r *TicketRepository) Create(ctx context.Context, ticket *models.Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

// GetByID возвращает заявку по идентификатору.
func (r *TicketRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.WithContext(ctx).First(&ticket, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

// Save сохраняет изменения заявки.
func (r *TicketRepository) Save(ctx context.Context, ticket *models.Ticket) error {
	return r.db.WithContext(ctx).Save(ticket).Error
}

// TicketListFilter описывает параметры фильтрации списка заявок.
type TicketListFilter struct {
	Status         string
	ContractorID   *uuid.UUID
	CleaningAreaID *uuid.UUID
	CreatedByOrgID *uuid.UUID
	DriverID       *uuid.UUID
}

// List возвращает список заявок с учётом заданных фильтров.
func (r *TicketRepository) List(ctx context.Context, filter TicketListFilter) ([]models.Ticket, error) {
	var tickets []models.Ticket

	query := r.db.WithContext(ctx).Model(&models.Ticket{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.ContractorID != nil {
		query = query.Where("contractor_id = ?", *filter.ContractorID)
	}
	if filter.CleaningAreaID != nil {
		query = query.Where("cleaning_area_id = ?", *filter.CleaningAreaID)
	}
	if filter.CreatedByOrgID != nil {
		query = query.Where("created_by_org_id = ?", *filter.CreatedByOrgID)
	}
	if filter.DriverID != nil {
		query = query.Joins("JOIN ticket_assignments ta ON ta.ticket_id = tickets.id").
			Where("ta.driver_id = ?", *filter.DriverID).
			Distinct()
	}

	if err := query.Order("created_at DESC").Find(&tickets).Error; err != nil {
		return nil, err
	}

	return tickets, nil
}
