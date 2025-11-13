package services

import (
	"fmt"

	"ticket-service/internal/models"
)

// CanTransition определяет, допустим ли переход статуса из current в next.
func CanTransition(current, next models.TicketStatus) bool {
	return ValidateStatusTransition(current, next) == nil
}

// ValidateStatusTransition проверяет допустимость перехода статусов.
func ValidateStatusTransition(oldStatus, newStatus models.TicketStatus) error {
	if oldStatus == newStatus {
		return nil
	}

	switch oldStatus {
	case models.TicketStatusPlanned:
		if newStatus == models.TicketStatusInProgress || newStatus == models.TicketStatusCancelled {
			return nil
		}
	case models.TicketStatusInProgress:
		if newStatus == models.TicketStatusCompleted || newStatus == models.TicketStatusCancelled {
			return nil
		}
	case models.TicketStatusCompleted:
		if newStatus == models.TicketStatusClosed {
			return nil
		}
	case models.TicketStatusCancelled, models.TicketStatusClosed:
		if newStatus == oldStatus {
			return nil
		}
	}

	return fmt.Errorf("status transition from %s to %s is not allowed", oldStatus, newStatus)
}
