package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"ticket-service/internal/database"
	"ticket-service/internal/models"
	"ticket-service/internal/repository"
	"ticket-service/internal/services"
)

// UpdateTicketStatusRequest описывает payload обновления статуса заявки.
type UpdateTicketStatusRequest struct {
	Status    string   `json:"status,omitempty"`
	Action    string   `json:"action,omitempty"`
	PhotoURL  *string  `json:"photo_url,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

// CreateTicketRequest — тело запроса для создания новой заявки.
type CreateTicketRequest struct {
	CleaningAreaID string  `json:"cleaning_area_id" binding:"required"`
	ContractorID   string  `json:"contractor_id" binding:"required"`
	PlannedStartAt string  `json:"planned_start_at" binding:"required"`
	PlannedEndAt   string  `json:"planned_end_at" binding:"required"`
	Description    *string `json:"description,omitempty"`
}

func getUserRole(c *gin.Context) string {
	if value, exists := c.Get("role"); exists {
		if role, ok := value.(string); ok && role != "" {
			return role
		}
	}
	return strings.TrimSpace(c.GetHeader("X-User-Role"))
}

func extractUUID(value interface{}) *uuid.UUID {
	switch v := value.(type) {
	case uuid.UUID:
		id := v
		return &id
	case *uuid.UUID:
		if v != nil {
			id := *v
			return &id
		}
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			return nil
		}
		if parsed, err := uuid.Parse(trimmed); err == nil {
			return &parsed
		}
	}
	return nil
}

func getOrgID(c *gin.Context) *uuid.UUID {
	if value, exists := c.Get("orgID"); exists {
		if id := extractUUID(value); id != nil {
			return id
		}
	}

	header := strings.TrimSpace(c.GetHeader("X-Org-ID"))
	if header == "" {
		return nil
	}

	if parsed, err := uuid.Parse(header); err == nil {
		return &parsed
	}

	return nil
}

func getDriverID(c *gin.Context) *uuid.UUID {
	if value, exists := c.Get("driverID"); exists {
		if id := extractUUID(value); id != nil {
			return id
		}
	}

	header := strings.TrimSpace(c.GetHeader("X-Driver-ID"))
	if header == "" {
		return nil
	}

	if parsed, err := uuid.Parse(header); err == nil {
		return &parsed
	}

	return nil
}

func parseUUIDQueryParam(c *gin.Context, name string) (*uuid.UUID, error) {
	val := strings.TrimSpace(c.Query(name))
	if val == "" {
		return nil, nil
	}

	parsed, err := uuid.Parse(val)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func parseTicketStatus(value string) (models.TicketStatus, bool) {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case string(models.TicketStatusPlanned):
		return models.TicketStatusPlanned, true
	case string(models.TicketStatusInProgress):
		return models.TicketStatusInProgress, true
	case string(models.TicketStatusCompleted):
		return models.TicketStatusCompleted, true
	case string(models.TicketStatusClosed):
		return models.TicketStatusClosed, true
	case string(models.TicketStatusCancelled):
		return models.TicketStatusCancelled, true
	default:
		return "", false
	}
}

// CreateTicketHandler обрабатывает запрос на создание новой заявки.
func CreateTicketHandler(c *gin.Context) {
	role := strings.ToUpper(getUserRole(c))
	if role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing role"})
		return
	}

	if role != "TOO_ADMIN" && role != "AKIMAT_ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only TOO or AKIMAT admin can create tickets"})
		return
	}

	orgIDHeader := strings.TrimSpace(c.GetHeader("X-Org-ID"))
	if orgIDHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing X-Org-ID header"})
		return
	}

	orgID, err := uuid.Parse(orgIDHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid X-Org-ID"})
		return
	}

	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	cleaningAreaID, err := uuid.Parse(strings.TrimSpace(req.CleaningAreaID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cleaning_area_id"})
		return
	}

	contractorID, err := uuid.Parse(strings.TrimSpace(req.ContractorID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contractor_id"})
		return
	}

	plannedStart, err := time.Parse(time.RFC3339, strings.TrimSpace(req.PlannedStartAt))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid planned_start_at"})
		return
	}

	plannedEnd, err := time.Parse(time.RFC3339, strings.TrimSpace(req.PlannedEndAt))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid planned_end_at"})
		return
	}

	ticket := &models.Ticket{
		CleaningAreaID: cleaningAreaID,
		ContractorID:   contractorID,
		CreatedByOrgID: orgID,
		Status:         models.TicketStatusPlanned,
		PlannedStartAt: plannedStart,
		PlannedEndAt:   plannedEnd,
	}

	if req.Description != nil {
		ticket.Description = strings.TrimSpace(*req.Description)
	}

	repo := repository.NewTicketRepository(database.GetDB())

	if err := repo.Create(c.Request.Context(), ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create ticket"})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

// GetTicketHandler возвращает заявку по идентификатору.
func GetTicketHandler(c *gin.Context) {
	idParam := c.Param("id")
	ticketID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	repo := repository.NewTicketRepository(database.GetDB())

	ticket, err := repo.GetByID(c.Request.Context(), ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load ticket"})
		return
	}

	// TODO: ограничить доступ к заявкам на основе ролей и принадлежности к организации.

	c.JSON(http.StatusOK, ticket)
}

// ListTicketsHandler возвращает список заявок с простейшими фильтрами.
func ListTicketsHandler(c *gin.Context) {
	role := strings.ToUpper(getUserRole(c))
	if role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing role"})
		return
	}

	statusParam := strings.TrimSpace(c.Query("status"))
	var statusFilter string
	if statusParam != "" {
		parsedStatus, ok := parseTicketStatus(statusParam)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown status"})
			return
		}
		statusFilter = string(parsedStatus)
	}

	contractorID, err := parseUUIDQueryParam(c, "contractor_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contractor_id"})
		return
	}

	cleaningAreaID, err := parseUUIDQueryParam(c, "cleaning_area_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cleaning_area_id"})
		return
	}

	filter := repository.TicketListFilter{
		Status:         statusFilter,
		ContractorID:   contractorID,
		CleaningAreaID: cleaningAreaID,
	}

	orgID := getOrgID(c)
	driverID := getDriverID(c)

	switch role {
	case "AKIMAT_ADMIN":
		// полный доступ
	case "TOO_ADMIN":
		if orgID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing organization id"})
			return
		}
		filter.CreatedByOrgID = orgID
	case "CONTRACTOR_ADMIN":
		if orgID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing organization id"})
			return
		}
		filter.ContractorID = orgID
	case "DRIVER":
		if driverID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "missing driver id"})
			return
		}
		filter.DriverID = driverID
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "role is not allowed"})
		return
	}

	repo := repository.NewTicketRepository(database.GetDB())

	tickets, err := repo.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// UpdateTicketStatusHandler обрабатывает запрос изменения статуса заявки.
func UpdateTicketStatusHandler(c *gin.Context) {
	idParam := c.Param("id")
	ticketID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	var req UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	targetStatus, err := resolveTargetStatus(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := strings.ToUpper(getUserRole(c))
	if !isRoleAllowed(role, targetStatus) {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	repo := repository.NewTicketRepository(database.GetDB())

	ticket, err := repo.GetByID(c.Request.Context(), ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load ticket"})
		return
	}

	if err := services.ValidateStatusTransition(ticket.Status, targetStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          err.Error(),
			"current_status": ticket.Status,
			"target_status":  targetStatus,
		})
		return
	}

	if ticket.Status == models.TicketStatusPlanned && targetStatus == models.TicketStatusInProgress && ticket.FactStartAt == nil {
		start := time.Now()
		ticket.FactStartAt = &start
	}

	if ticket.Status == models.TicketStatusInProgress && targetStatus == models.TicketStatusCompleted && ticket.FactEndAt == nil {
		end := time.Now()
		ticket.FactEndAt = &end
	}

	ticket.Status = targetStatus

	if req.PhotoURL != nil {
		ticket.PhotoURL = req.PhotoURL
	}
	if req.Latitude != nil {
		ticket.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		ticket.Longitude = req.Longitude
	}

	if err := repo.Save(c.Request.Context(), ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update ticket"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func resolveTargetStatus(req *UpdateTicketStatusRequest) (models.TicketStatus, error) {
	var zero models.TicketStatus

	status := strings.TrimSpace(req.Status)
	if status != "" {
		if parsed, ok := parseTicketStatus(status); ok {
			return parsed, nil
		}
		return zero, errors.New("unknown status")
	}

	action := strings.ToLower(strings.TrimSpace(req.Action))
	switch action {
	case "mark_in_progress":
		return models.TicketStatusInProgress, nil
	case "mark_completed":
		return models.TicketStatusCompleted, nil
	case "mark_closed":
		return models.TicketStatusClosed, nil
	case "cancel":
		return models.TicketStatusCancelled, nil
	case "":
		return zero, errors.New("status or action must be provided")
	default:
		return zero, errors.New("unknown action")
	}
}

func isRoleAllowed(role string, targetStatus models.TicketStatus) bool {
	switch targetStatus {
	case models.TicketStatusInProgress, models.TicketStatusCompleted:
		return role == "DRIVER" || role == "CONTRACTOR_ADMIN"
	case models.TicketStatusClosed, models.TicketStatusCancelled:
		return role == "TOO_ADMIN" || role == "AKIMAT_ADMIN"
	case models.TicketStatusPlanned:
		return role != ""
	default:
		return false
	}
}

// RegisterTicketRoutes регистрирует HTTP-маршруты для работы с заявками.
func RegisterTicketRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	tickets := api.Group("/tickets")
	{
		tickets.POST("", CreateTicketHandler)
		tickets.GET("", ListTicketsHandler)
		tickets.GET("/:id", GetTicketHandler)
		tickets.PATCH("/:id/status", UpdateTicketStatusHandler)
	}

	akimat := api.Group("/akimat")
	akimat.GET("/tickets", ListTicketsHandler)

	too := api.Group("/too")
	too.GET("/tickets", ListTicketsHandler)

	contractor := api.Group("/contractor")
	contractor.GET("/tickets", ListTicketsHandler)

	driver := api.Group("/driver")
	driver.GET("/tickets", ListTicketsHandler)
}
