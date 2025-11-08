package handlers

import (
	"net/http"
	"strconv"

	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/connect-univyn/connect_server/internal/service/analytics"
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnalyticsHandler struct {
	analyticsService *analytics.Service
}

func NewAnalyticsHandler(analyticsService *analytics.Service) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// ============================================================================
// Content Moderation & Reporting Handlers
// ============================================================================

// CreateReport godoc
// @Summary Submit a content report
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /api/analytics/reports [post]
func (h *AnalyticsHandler) CreateReport(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	var req analytics.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid user ID"))
		return
	}
	req.ReporterID = userID

	report, err := h.analyticsService.CreateReport(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(report))
}

// GetReport godoc
// @Summary Get report by ID
// @Tags analytics
// @Produce json
// @Param id path string true "Report ID"
// @Router /api/analytics/reports/:id [get]
func (h *AnalyticsHandler) GetReport(c *gin.Context) {
	reportID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid report ID"))
		return
	}

	report, err := h.analyticsService.GetReport(c.Request.Context(), reportID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(report))
}

// GetReportsByContent godoc
// @Summary Get reports for specific content
// @Tags analytics
// @Produce json
// @Param content_type query string true "Content type"
// @Param content_id query string true "Content ID"
// @Router /api/analytics/reports/by-content [get]
func (h *AnalyticsHandler) GetReportsByContent(c *gin.Context) {
	contentType := c.Query("content_type")
	if contentType == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "content_type is required"))
		return
	}

	contentID, err := uuid.Parse(c.Query("content_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid content_id"))
		return
	}

	reports, err := h.analyticsService.GetReportsByContent(c.Request.Context(), contentType, contentID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(reports))
}

// GetModerationQueue godoc
// @Summary Get moderation queue
// @Tags analytics
// @Produce json
// @Param space_id query string true "Space ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Router /api/analytics/moderation/queue [get]
func (h *AnalyticsHandler) GetModerationQueue(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	page := int32(1)
	if pageStr := c.Query("page"); pageStr != "" {
		p, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil || p < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid page number"))
			return
		}
		page = int32(p)
	}

	limit := int32(20)
	if limitStr := c.Query("limit"); limitStr != "" {
		l, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || l < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid limit"))
			return
		}
		limit = int32(l)
	}

	reports, err := h.analyticsService.GetModerationQueue(c.Request.Context(), spaceID, page, limit)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(reports))
}

// GetPendingReports godoc
// @Summary Get pending reports
// @Tags analytics
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/analytics/reports/pending [get]
func (h *AnalyticsHandler) GetPendingReports(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	reports, err := h.analyticsService.GetPendingReports(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(reports))
}

// UpdateReport godoc
// @Summary Update report (moderate)
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Report ID"
// @Router /api/analytics/reports/:id [put]
func (h *AnalyticsHandler) UpdateReport(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	reviewerID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid user ID"))
		return
	}

	reportID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid report ID"))
		return
	}

	var req analytics.UpdateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	report, err := h.analyticsService.UpdateReport(c.Request.Context(), reportID, reviewerID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(report))
}

// GetContentModerationStats godoc
// @Summary Get moderation statistics
// @Tags analytics
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/analytics/moderation/stats [get]
func (h *AnalyticsHandler) GetContentModerationStats(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	stats, err := h.analyticsService.GetContentModerationStats(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(stats))
}

// ============================================================================
// System & Space Metrics Handlers
// ============================================================================

// GetSystemMetrics godoc
// @Summary Get current system metrics
// @Tags analytics
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/analytics/metrics/system [get]
func (h *AnalyticsHandler) GetSystemMetrics(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	metrics, err := h.analyticsService.GetSystemMetrics(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(metrics))
}

// GetSpaceStats godoc
// @Summary Get space statistics
// @Tags analytics
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/analytics/metrics/space [get]
func (h *AnalyticsHandler) GetSpaceStats(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	stats, err := h.analyticsService.GetSpaceStats(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(stats))
}

// ============================================================================
// Engagement & Activity Handlers
// ============================================================================

// GetEngagementMetrics godoc
// @Summary Get engagement metrics over time
// @Tags analytics
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/analytics/engagement/metrics [get]
func (h *AnalyticsHandler) GetEngagementMetrics(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	metrics, err := h.analyticsService.GetEngagementMetrics(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(metrics))
}

// GetUserActivityStats godoc
// @Summary Get user activity statistics
// @Tags analytics
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/analytics/activity/stats [get]
func (h *AnalyticsHandler) GetUserActivityStats(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	stats, err := h.analyticsService.GetUserActivityStats(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(stats))
}

// GetUserGrowth godoc
// @Summary Get user growth over time
// @Tags analytics
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/analytics/users/growth [get]
func (h *AnalyticsHandler) GetUserGrowth(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	growth, err := h.analyticsService.GetUserGrowth(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(growth))
}

// GetUserEngagementRanking godoc
// @Summary Get top users by engagement
// @Tags analytics
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/analytics/users/ranking [get]
func (h *AnalyticsHandler) GetUserEngagementRanking(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	ranking, err := h.analyticsService.GetUserEngagementRanking(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(ranking))
}

// ============================================================================
// Top Content Handlers
// ============================================================================

// GetTopPosts godoc
// @Summary Get most engaging posts
// @Tags analytics
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/analytics/top/posts [get]
func (h *AnalyticsHandler) GetTopPosts(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	posts, err := h.analyticsService.GetTopPosts(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(posts))
}

// GetTopCommunities godoc
// @Summary Get most engaging communities
// @Tags analytics
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/analytics/top/communities [get]
func (h *AnalyticsHandler) GetTopCommunities(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	communities, err := h.analyticsService.GetTopCommunities(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(communities))
}

// GetTopGroups godoc
// @Summary Get most active groups
// @Tags analytics
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/analytics/top/groups [get]
func (h *AnalyticsHandler) GetTopGroups(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	groups, err := h.analyticsService.GetTopGroups(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(groups))
}

// ============================================================================
// Mentorship Analytics Handlers
// ============================================================================

// GetMentoringStats godoc
// @Summary Get mentoring statistics
// @Tags analytics
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/analytics/mentorship/mentoring [get]
func (h *AnalyticsHandler) GetMentoringStats(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	stats, err := h.analyticsService.GetMentoringStats(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(stats))
}

// GetTutoringStats godoc
// @Summary Get tutoring statistics
// @Tags analytics
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/analytics/mentorship/tutoring [get]
func (h *AnalyticsHandler) GetTutoringStats(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	stats, err := h.analyticsService.GetTutoringStats(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(stats))
}

// GetPopularIndustries godoc
// @Summary Get popular mentoring industries
// @Tags analytics
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/analytics/mentorship/industries [get]
func (h *AnalyticsHandler) GetPopularIndustries(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	industries, err := h.analyticsService.GetPopularIndustries(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(industries))
}

// GetPopularSubjects godoc
// @Summary Get popular tutoring subjects
// @Tags analytics
// @Produce json
// @Param space_id query string true "Space ID"
// @Router /api/analytics/mentorship/subjects [get]
func (h *AnalyticsHandler) GetPopularSubjects(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid space ID"))
		return
	}

	subjects, err := h.analyticsService.GetPopularSubjects(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(subjects))
}
