package handlers

import (
	"net/http"
	"strconv"

	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/connect-univyn/connect-server/internal/service/analytics"
	"github.com/connect-univyn/connect-server/internal/util"
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
