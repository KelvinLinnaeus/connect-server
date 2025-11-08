package handlers

import (
	"net/http"
	"strconv"

	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/connect-univyn/connect_server/internal/service/announcements"
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnnouncementHandler struct {
	announcementService *announcements.Service
}

func NewAnnouncementHandler(announcementService *announcements.Service) *AnnouncementHandler {
	return &AnnouncementHandler{
		announcementService: announcementService,
	}
}

// CreateAnnouncement godoc
// @Summary Create announcement
// @Tags announcements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /api/announcements [post]
func (h *AnnouncementHandler) CreateAnnouncement(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	var req announcements.CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid user ID"))
		return
	}
	req.AuthorID = userID

	announcement, err := h.announcementService.CreateAnnouncement(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(announcement))
}

// GetAnnouncement godoc
// @Summary Get announcement by ID
// @Tags announcements
// @Produce json
// @Param id path string true "Announcement ID"
// @Router /api/announcements/:id [get]
func (h *AnnouncementHandler) GetAnnouncement(c *gin.Context) {
	announcementID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid announcement ID"))
		return
	}

	announcement, err := h.announcementService.GetAnnouncementByID(c.Request.Context(), announcementID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(announcement))
}

// ListAnnouncements godoc
// @Summary List announcements
// @Tags announcements
// @Produce json
// @Param space_id query string true "Space ID"
// @Param target_audience query string false "Target audience filter (comma-separated)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Router /api/announcements [get]
func (h *AnnouncementHandler) ListAnnouncements(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	params := announcements.ListAnnouncementsParams{
		SpaceID: spaceID,
		Page:    1,
		Limit:   20,
	}

	if targetAudience := c.Query("target_audience"); targetAudience != "" {
		// Parse comma-separated target audience
		// For simplicity, we'll accept a single value for now
		params.TargetAudience = []string{targetAudience}
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		params.Page = int32(page)
	}

	if limit, err := strconv.Atoi(c.Query("limit")); err == nil && limit > 0 {
		params.Limit = int32(limit)
	}

	announcementsList, err := h.announcementService.ListAnnouncements(c.Request.Context(), params)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(announcementsList))
}

// UpdateAnnouncement godoc
// @Summary Update announcement
// @Tags announcements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Announcement ID"
// @Router /api/announcements/:id [put]
func (h *AnnouncementHandler) UpdateAnnouncement(c *gin.Context) {
	announcementID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid announcement ID"))
		return
	}

	var req announcements.UpdateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	announcement, err := h.announcementService.UpdateAnnouncement(c.Request.Context(), announcementID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(announcement))
}

// UpdateAnnouncementStatus godoc
// @Summary Update announcement status
// @Tags announcements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Announcement ID"
// @Router /api/announcements/:id/status [put]
func (h *AnnouncementHandler) UpdateAnnouncementStatus(c *gin.Context) {
	announcementID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid announcement ID"))
		return
	}

	var req announcements.UpdateAnnouncementStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	announcement, err := h.announcementService.UpdateAnnouncementStatus(c.Request.Context(), announcementID, req.Status)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(announcement))
}
