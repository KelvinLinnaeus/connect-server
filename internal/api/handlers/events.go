package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/connect-univyn/connect_server/internal/service/events"
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventHandler struct {
	eventService *events.Service
}

func NewEventHandler(eventService *events.Service) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

// CreateEvent godoc
// @Summary Create event
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /api/events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	var req events.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid user ID"))
		return
	}
	req.OrganizerID = userID

	event, err := h.eventService.CreateEvent(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(event))
}

// GetEvent godoc
// @Summary Get event by ID
// @Tags events
// @Produce json
// @Param id path string true "Event ID"
// @Router /api/events/:id [get]
func (h *EventHandler) GetEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid event ID"))
		return
	}

	// Try to get user ID from auth (optional for public events)
	userID := uuid.Nil
	if payload, exists := c.Get("authorization_payload"); exists {
		authPayload := payload.(*auth.Payload)
		if parsed, err := uuid.Parse(authPayload.UserID); err == nil {
			userID = parsed
		}
	}

	event, err := h.eventService.GetEventByID(c.Request.Context(), eventID, userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(event))
}

// ListEvents godoc
// @Summary List events
// @Tags events
// @Produce json
// @Param space_id query string true "Space ID"
// @Param category query string false "Category filter"
// @Param start_date query string false "Start date filter (ISO 8601)"
// @Param sort query string false "Sort by: upcoming, popular, recent"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Router /api/events [get]
func (h *EventHandler) ListEvents(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	userID := uuid.Nil
	if payload, exists := c.Get("authorization_payload"); exists {
		authPayload := payload.(*auth.Payload)
		if parsed, err := uuid.Parse(authPayload.UserID); err == nil {
			userID = parsed
		}
	}

	params := events.ListEventsParams{
		SpaceID: spaceID,
		UserID:  userID,
		Page:    1,
		Limit:   20,
	}

	if category := c.Query("category"); category != "" {
		params.Category = &category
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			params.StartDate = &startDate
		}
	}

	if sort := c.Query("sort"); sort != "" {
		params.Sort = &sort
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		params.Page = int32(page)
	}

	if limit, err := strconv.Atoi(c.Query("limit")); err == nil && limit > 0 {
		params.Limit = int32(limit)
	}

	eventsList, err := h.eventService.ListEvents(c.Request.Context(), params)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(eventsList))
}

// GetUpcomingEvents godoc
// @Summary Get upcoming events (next 7 days)
// @Tags events
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/events/upcoming [get]
func (h *EventHandler) GetUpcomingEvents(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	eventsList, err := h.eventService.GetUpcomingEvents(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(eventsList))
}

// GetUserEvents godoc
// @Summary Get user's registered events
// @Tags events
// @Produce json
// @Security BearerAuth
// @Param space_id query string true "Space ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Router /api/users/events [get]
func (h *EventHandler) GetUserEvents(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid user ID"))
		return
	}

	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	page := int32(1)
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = int32(p)
	}

	limit := int32(20)
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = int32(l)
	}

	eventsList, err := h.eventService.GetUserEvents(c.Request.Context(), userID, spaceID, page, limit)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(eventsList))
}

// SearchEvents godoc
// @Summary Search events
// @Tags events
// @Produce json
// @Param space_id query string true "Space ID"
// @Param q query string true "Search query"
// @Router /api/events/search [get]
func (h *EventHandler) SearchEvents(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Search query required"))
		return
	}

	userID := uuid.Nil
	if payload, exists := c.Get("authorization_payload"); exists {
		authPayload := payload.(*auth.Payload)
		if parsed, err := uuid.Parse(authPayload.UserID); err == nil {
			userID = parsed
		}
	}

	params := events.SearchEventsParams{
		SpaceID: spaceID,
		UserID:  userID,
		Query:   query,
	}

	eventsList, err := h.eventService.SearchEvents(c.Request.Context(), params)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(eventsList))
}

// GetEventCategories godoc
// @Summary Get event categories
// @Tags events
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/events/categories [get]
func (h *EventHandler) GetEventCategories(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	categories, err := h.eventService.GetEventCategories(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(categories))
}

// RegisterForEvent godoc
// @Summary Register for event
// @Tags events
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Router /api/events/:id/register [post]
func (h *EventHandler) RegisterForEvent(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid user ID"))
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid event ID"))
		return
	}

	registration, err := h.eventService.RegisterForEvent(c.Request.Context(), eventID, userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(registration))
}

// UnregisterFromEvent godoc
// @Summary Unregister from event
// @Tags events
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Router /api/events/:id/unregister [post]
func (h *EventHandler) UnregisterFromEvent(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid user ID"))
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid event ID"))
		return
	}

	err = h.eventService.UnregisterFromEvent(c.Request.Context(), eventID, userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Successfully unregistered from event"}))
}

// GetEventAttendees godoc
// @Summary Get event attendees
// @Tags events
// @Produce json
// @Param id path string true "Event ID"
// @Router /api/events/:id/attendees [get]
func (h *EventHandler) GetEventAttendees(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid event ID"))
		return
	}

	attendees, err := h.eventService.GetEventAttendees(c.Request.Context(), eventID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(attendees))
}

// MarkEventAttendance godoc
// @Summary Mark user as attended
// @Tags events
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Param user_id path string true "User ID"
// @Router /api/events/:id/attendance/:user_id [post]
func (h *EventHandler) MarkEventAttendance(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid event ID"))
		return
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid user ID"))
		return
	}

	err = h.eventService.MarkEventAttendance(c.Request.Context(), eventID, userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Attendance marked successfully"}))
}

// AddEventCoOrganizer godoc
// @Summary Add co-organizer to event
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Router /api/events/:id/co-organizers [post]
func (h *EventHandler) AddEventCoOrganizer(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid event ID"))
		return
	}

	var req events.AddCoOrganizerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	coOrganizer, err := h.eventService.AddEventCoOrganizer(c.Request.Context(), eventID, req.UserID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(coOrganizer))
}

// GetEventCoOrganizers godoc
// @Summary Get event co-organizers
// @Tags events
// @Produce json
// @Param id path string true "Event ID"
// @Router /api/events/:id/co-organizers [get]
func (h *EventHandler) GetEventCoOrganizers(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid event ID"))
		return
	}

	coOrganizers, err := h.eventService.GetEventCoOrganizers(c.Request.Context(), eventID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(coOrganizers))
}

// RemoveEventCoOrganizer godoc
// @Summary Remove co-organizer from event
// @Tags events
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Param user_id path string true "User ID"
// @Router /api/events/:id/co-organizers/:user_id [delete]
func (h *EventHandler) RemoveEventCoOrganizer(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid event ID"))
		return
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid user ID"))
		return
	}

	err = h.eventService.RemoveEventCoOrganizer(c.Request.Context(), eventID, userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Co-organizer removed successfully"}))
}

// UpdateEvent godoc
// @Summary Update event
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Router /api/events/:id [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid event ID"))
		return
	}

	var req events.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	event, err := h.eventService.UpdateEvent(c.Request.Context(), eventID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(event))
}

// UpdateEventStatus godoc
// @Summary Update event status
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Router /api/events/:id/status [put]
func (h *EventHandler) UpdateEventStatus(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid event ID"))
		return
	}

	var req events.UpdateEventStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	event, err := h.eventService.UpdateEventStatus(c.Request.Context(), eventID, req.Status)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(event))
}
