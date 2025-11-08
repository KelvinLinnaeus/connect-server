package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/connect-univyn/connect_server/internal/service/admin"
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	adminService *admin.Service
}

func NewAdminHandler(adminService *admin.Service) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// PUT /api/admin/users/:id/suspend
func (h *AdminHandler) SuspendUser(c *gin.Context) {
	// Get authenticated admin user
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "No auth payload"))
		return
	}
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	// Parse user ID
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid user ID"))
		return
	}

	// Parse request body
	var req struct {
		Reason       string `json:"reason" binding:"required"`
		Notes        string `json:"notes"`
		DurationDays int    `json:"duration_days"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	// Suspend user
	err = h.adminService.SuspendUser(c.Request.Context(), admin.SuspendUserRequest{
		UserID:       userID,
		SuspendedBy:  adminUserID,
		Reason:       req.Reason,
		Notes:        req.Notes,
		DurationDays: req.DurationDays,
		IsPermanent:  req.DurationDays == 0,
	})
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "User suspended successfully",
	}))
}

// PUT /api/admin/users/:id/unsuspend
func (h *AdminHandler) UnsuspendUser(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid user ID"))
		return
	}

	err = h.adminService.UnsuspendUser(c.Request.Context(), userID, adminUserID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "User unsuspended successfully",
	}))
}

// PUT /api/admin/users/:id/ban
func (h *AdminHandler) BanUser(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid user ID"))
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	err = h.adminService.BanUser(c.Request.Context(), userID, adminUserID, req.Reason)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "User banned successfully",
	}))
}

// GET /api/admin/reports
func (h *AdminHandler) GetReports(c *gin.Context) {
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	contentType := c.Query("content_type")

	offset := (page - 1) * limit

	reports, err := h.adminService.GetContentReports(c.Request.Context(), admin.GetReportsRequest{
		SpaceID:     spaceID,
		Status:      status,
		ContentType: contentType,
		Limit:       int32(limit),
		Offset:      int32(offset),
	})
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(map[string]interface{}{
		"reports": reports,
		"page":    page,
		"limit":   limit,
	}))
}

// GET /api/admin/spaces/:id/activities
func (h *AdminHandler) GetSpaceActivities(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid space ID"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	activityType := c.Query("activity_type")
	since := c.Query("since")

	var sinceTime time.Time
	if since != "" {
		sinceTime, _ = time.Parse(time.RFC3339, since)
	}

	offset := (page - 1) * limit

	activities, err := h.adminService.GetSpaceActivities(
		c.Request.Context(),
		spaceID,
		activityType,
		sinceTime,
		int32(limit),
		int32(offset),
	)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(map[string]interface{}{
		"activities": activities,
		"page":       page,
		"limit":      limit,
	}))
}

// GET /api/admin/dashboard/stats
func (h *AdminHandler) GetDashboardStats(c *gin.Context) {
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID"))
		return
	}

	stats, err := h.adminService.GetDashboardStats(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(stats))
}

// GET /api/admin/users
func (h *AdminHandler) GetUsers(c *gin.Context) {
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	users, total, err := h.adminService.GetUsers(c.Request.Context(), spaceID, int32(limit), int32(offset))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(map[string]interface{}{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	}))
}

// DELETE /api/admin/users/:id
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid user ID"))
		return
	}

	err = h.adminService.DeleteUser(c.Request.Context(), userID, adminUserID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "User deleted successfully",
	}))
}

// GET /api/admin/applications/tutors
func (h *AdminHandler) GetTutorApplications(c *gin.Context) {
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	applications, err := h.adminService.GetTutorApplications(c.Request.Context(), spaceID, int32(limit), int32(offset))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(map[string]interface{}{
		"applications": applications,
		"page":         page,
		"limit":        limit,
	}))
}

// PUT /api/admin/applications/tutors/:id/approve
func (h *AdminHandler) ApproveTutorApplication(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid application ID"))
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	err = h.adminService.ApproveTutorApplication(c.Request.Context(), appID, adminUserID, req.Notes)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Tutor application approved successfully",
	}))
}

// PUT /api/admin/applications/tutors/:id/reject
func (h *AdminHandler) RejectTutorApplication(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid application ID"))
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	err = h.adminService.RejectTutorApplication(c.Request.Context(), appID, adminUserID, req.Notes)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Tutor application rejected successfully",
	}))
}

// GET /api/admin/applications/mentors
func (h *AdminHandler) GetMentorApplications(c *gin.Context) {
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	applications, err := h.adminService.GetMentorApplications(c.Request.Context(), spaceID, int32(limit), int32(offset))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(map[string]interface{}{
		"applications": applications,
		"page":         page,
		"limit":        limit,
	}))
}

// PUT /api/admin/applications/mentors/:id/approve
func (h *AdminHandler) ApproveMentorApplication(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid application ID"))
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	err = h.adminService.ApproveMentorApplication(c.Request.Context(), appID, adminUserID, req.Notes)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Mentor application approved successfully",
	}))
}

// PUT /api/admin/applications/mentors/:id/reject
func (h *AdminHandler) RejectMentorApplication(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid application ID"))
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	err = h.adminService.RejectMentorApplication(c.Request.Context(), appID, adminUserID, req.Notes)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Mentor application rejected successfully",
	}))
}

// PUT /api/admin/reports/:id/resolve
func (h *AdminHandler) ResolveReport(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	reportID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid report ID"))
		return
	}

	var req struct {
		Action string `json:"action"`
		Notes  string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	err = h.adminService.ResolveReport(c.Request.Context(), reportID, adminUserID, req.Action, req.Notes)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Report resolved successfully",
	}))
}

// PUT /api/admin/reports/:id/escalate
func (h *AdminHandler) EscalateReport(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	reportID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid report ID"))
		return
	}

	err = h.adminService.EscalateReport(c.Request.Context(), reportID, adminUserID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Report escalated successfully",
	}))
}

// GET /api/admin/groups
func (h *AdminHandler) GetGroups(c *gin.Context) {
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	offset := (page - 1) * limit

	groups, total, err := h.adminService.GetGroups(c.Request.Context(), spaceID, status, int32(limit), int32(offset))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(map[string]interface{}{
		"groups": groups,
		"total":  total,
		"page":   page,
		"limit":  limit,
	}))
}

// PUT /api/admin/groups/:id/approve
func (h *AdminHandler) ApproveGroup(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID"))
		return
	}

	err = h.adminService.ApproveGroup(c.Request.Context(), groupID, adminUserID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Group approved successfully",
	}))
}

// PUT /api/admin/groups/:id/reject
func (h *AdminHandler) RejectGroup(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID"))
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	err = h.adminService.RejectGroup(c.Request.Context(), groupID, adminUserID, req.Reason)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Group rejected successfully",
	}))
}

// DELETE /api/admin/groups/:id
func (h *AdminHandler) DeleteGroup(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID"))
		return
	}

	err = h.adminService.DeleteGroup(c.Request.Context(), groupID, adminUserID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Group deleted successfully",
	}))
}

// System Settings Handlers

// GET /api/admin/settings
func (h *AdminHandler) GetSettings(c *gin.Context) {
	settings, err := h.adminService.GetAllSettings(c.Request.Context())
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"settings": settings,
	}))
}

// PUT /api/admin/settings/:key
func (h *AdminHandler) UpdateSetting(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	key := c.Param("key")

	var req struct {
		Value       map[string]interface{} `json:"value" binding:"required"`
		Description string                 `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	// Convert value to JSON
	valueJSON, err := json.Marshal(req.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_value", "Invalid value format"))
		return
	}

	setting, err := h.adminService.UpdateSetting(c.Request.Context(), key, valueJSON, req.Description, adminUserID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Setting updated successfully",
		"setting": setting,
	}))
}

// Analytics Handlers

// GET /api/admin/analytics/user-growth
func (h *AdminHandler) GetUserGrowth(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID"))
		return
	}

	// Default to last 6 months
	since := time.Now().AddDate(0, -6, 0)
	if sinceParam := c.Query("since"); sinceParam != "" {
		parsedSince, err := time.Parse(time.RFC3339, sinceParam)
		if err == nil {
			since = parsedSince
		}
	}

	data, err := h.adminService.GetUserGrowth(c.Request.Context(), spaceID, since)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"data": data,
	}))
}

// GET /api/admin/analytics/engagement
func (h *AdminHandler) GetEngagementMetrics(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID"))
		return
	}

	since := time.Now().AddDate(0, -1, 0) // Last month by default
	if sinceParam := c.Query("since"); sinceParam != "" {
		parsedSince, err := time.Parse(time.RFC3339, sinceParam)
		if err == nil {
			since = parsedSince
		}
	}

	data, err := h.adminService.GetContentGrowth(c.Request.Context(), spaceID, since)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"data": data,
	}))
}

// GET /api/admin/analytics/activity
func (h *AdminHandler) GetActivityAnalytics(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID"))
		return
	}

	since := time.Now().AddDate(0, -1, 0) // Last month by default
	if sinceParam := c.Query("since"); sinceParam != "" {
		parsedSince, err := time.Parse(time.RFC3339, sinceParam)
		if err == nil {
			since = parsedSince
		}
	}

	stats, err := h.adminService.GetActivityStats(c.Request.Context(), spaceID, since)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(stats))
}

// Admin Management Handlers

// GET /api/admin/admins
func (h *AdminHandler) GetAdmins(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	admins, err := h.adminService.GetAllAdmins(c.Request.Context(), status, int32(limit), int32(offset))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"admins": admins,
		"page":   page,
		"limit":  limit,
	}))
}

// PUT /api/admin/admins/:id/role
func (h *AdminHandler) UpdateAdminRole(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid user ID"))
		return
	}

	var req struct {
		Roles []string `json:"roles" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	err = h.adminService.UpdateAdminRole(c.Request.Context(), userID, adminUserID, req.Roles)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Admin role updated successfully",
	}))
}

// PUT /api/admin/admins/:id/status
func (h *AdminHandler) UpdateAdminStatus(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid user ID"))
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	err = h.adminService.UpdateAdminStatus(c.Request.Context(), userID, adminUserID, req.Status)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Admin status updated successfully",
	}))
}

// Notification Handlers

// GET /api/admin/notifications
func (h *AdminHandler) GetNotifications(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	userID, _ := uuid.Parse(authPayload.UserID)

	typeFilter := c.Query("type")
	priority := c.Query("priority")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var isRead *bool
	if c.Query("is_read") != "" {
		val := c.Query("is_read") == "true"
		isRead = &val
	}

	notifications, err := h.adminService.GetNotifications(c.Request.Context(), userID, typeFilter, priority, isRead, int32(limit), int32(offset))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"notifications": notifications,
		"page":          page,
		"limit":         limit,
	}))
}

// PUT /api/admin/notifications/:id/read
func (h *AdminHandler) MarkNotificationRead(c *gin.Context) {
	notificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid notification ID"))
		return
	}

	err = h.adminService.MarkNotificationAsRead(c.Request.Context(), notificationID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Notification marked as read",
	}))
}

// DELETE /api/admin/notifications/:id
func (h *AdminHandler) DeleteNotification(c *gin.Context) {
	notificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid notification ID"))
		return
	}

	err = h.adminService.DeleteNotification(c.Request.Context(), notificationID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Notification deleted successfully",
	}))
}

// PUT /api/admin/notifications/read-all
func (h *AdminHandler) MarkAllNotificationsRead(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	userID, _ := uuid.Parse(authPayload.UserID)

	err := h.adminService.MarkAllNotificationsAsRead(c.Request.Context(), userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "All notifications marked as read",
	}))
}

// ==================== Communities Management ====================

// GET /api/admin/communities
func (h *AdminHandler) GetCommunities(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID"))
		return
	}

	

	category := c.Query("category")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	communities, err := h.adminService.GetAllCommunities(c.Request.Context(), spaceID, category, status, int32(limit), int32(offset))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"communities": communities,
		"page":        page,
		"limit":       limit,
	}))
}

// POST /api/admin/communities
func (h *AdminHandler) CreateCommunity(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	var req admin.CreateCommunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_request", err.Error()))
		return
	}

	community, err := h.adminService.CreateCommunity(c.Request.Context(), adminID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(gin.H{
		"community": community,
	}))
}

// PUT /api/admin/communities/:id
func (h *AdminHandler) UpdateCommunity(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid community ID"))
		return
	}

	var req admin.UpdateCommunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_request", err.Error()))
		return
	}

	community, err := h.adminService.UpdateCommunity(c.Request.Context(), communityID, adminID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"community": community,
	}))
}

// DELETE /api/admin/communities/:id
func (h *AdminHandler) DeleteCommunity(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid community ID"))
		return
	}

	err = h.adminService.DeleteCommunity(c.Request.Context(), communityID, adminID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Community deleted successfully",
	}))
}

// PUT /api/admin/communities/:id/status
func (h *AdminHandler) UpdateCommunityStatus(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid community ID"))
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_request", err.Error()))
		return
	}

	community, err := h.adminService.UpdateCommunityStatus(c.Request.Context(), communityID, adminID, req.Status)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"community": community,
	}))
}

// POST /api/admin/communities/:id/moderators
func (h *AdminHandler) AssignCommunityModerator(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid community ID"))
		return
	}

	var req struct {
		UserID      string   `json:"user_id" binding:"required"`
		Permissions []string `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_request", err.Error()))
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user_id", "Invalid user ID"))
		return
	}

	err = h.adminService.AssignCommunityModerator(c.Request.Context(), communityID, userID, adminID, req.Permissions)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Moderator assigned successfully",
	}))
}

// ==================== Announcements Management ====================

// GET /api/admin/announcements
func (h *AdminHandler) GetAnnouncements(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID"))
		return
	}

	status := c.Query("status")
	priority := c.Query("priority")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	announcements, err := h.adminService.GetAllAnnouncements(c.Request.Context(), spaceID, status, priority, int32(limit), int32(offset))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"announcements": announcements,
		"page":          page,
		"limit":         limit,
	}))
}

// POST /api/admin/announcements
func (h *AdminHandler) CreateAnnouncement(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	var req admin.CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_request", err.Error()))
		return
	}

	announcement, err := h.adminService.CreateAnnouncement(c.Request.Context(), adminID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(gin.H{
		"announcement": announcement,
	}))
}

// PUT /api/admin/announcements/:id
func (h *AdminHandler) UpdateAnnouncement(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	announcementID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid announcement ID"))
		return
	}

	var req admin.UpdateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_request", err.Error()))
		return
	}

	announcement, err := h.adminService.UpdateAnnouncement(c.Request.Context(), announcementID, adminID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"announcement": announcement,
	}))
}

// DELETE /api/admin/announcements/:id
func (h *AdminHandler) DeleteAnnouncement(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	announcementID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid announcement ID"))
		return
	}

	err = h.adminService.DeleteAnnouncement(c.Request.Context(), announcementID, adminID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Announcement deleted successfully",
	}))
}

// PUT /api/admin/announcements/:id/status
func (h *AdminHandler) UpdateAnnouncementStatus(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	announcementID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid announcement ID"))
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_request", err.Error()))
		return
	}

	announcement, err := h.adminService.UpdateAnnouncementStatus(c.Request.Context(), announcementID, adminID, req.Status)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"announcement": announcement,
	}))
}

// ==================== Events Management ====================

// GET /api/admin/events
func (h *AdminHandler) GetEvents(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID"))
		return
	}

	status := c.Query("status")
	category := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	events, err := h.adminService.GetAllEvents(c.Request.Context(), spaceID, status, category, int32(limit), int32(offset))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"events": events,
		"page":   page,
		"limit":  limit,
	}))
}

// POST /api/admin/events
func (h *AdminHandler) CreateEvent(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	var req admin.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_request", err.Error()))
		return
	}

	event, err := h.adminService.CreateEvent(c.Request.Context(), adminID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(gin.H{
		"event": event,
	}))
}

// PUT /api/admin/events/:id
func (h *AdminHandler) UpdateEvent(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid event ID"))
		return
	}

	var req admin.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_request", err.Error()))
		return
	}

	event, err := h.adminService.UpdateEvent(c.Request.Context(), eventID, adminID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"event": event,
	}))
}

// DELETE /api/admin/events/:id
func (h *AdminHandler) DeleteEvent(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid event ID"))
		return
	}

	err = h.adminService.DeleteEvent(c.Request.Context(), eventID, adminID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Event deleted successfully",
	}))
}

// PUT /api/admin/events/:id/status
func (h *AdminHandler) UpdateEventStatus(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid event ID"))
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_request", err.Error()))
		return
	}

	event, err := h.adminService.UpdateEventStatus(c.Request.Context(), eventID, adminID, req.Status)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"event": event,
	}))
}

// GET /api/admin/events/:id/registrations
func (h *AdminHandler) GetEventRegistrations(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid event ID"))
		return
	}

	registrations, err := h.adminService.GetEventRegistrations(c.Request.Context(), eventID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"registrations": registrations,
	}))
}

// ==================== User Management ====================

// POST /api/admin/users
func (h *AdminHandler) CreateUser(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	var req admin.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_request", err.Error()))
		return
	}

	user, err := h.adminService.CreateUser(c.Request.Context(), adminID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(gin.H{
		"user": user,
	}))
}

// PUT /api/admin/users/:id
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid user ID"))
		return
	}

	var req admin.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_request", err.Error()))
		return
	}

	user, err := h.adminService.UpdateUser(c.Request.Context(), userID, adminID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"user": user,
	}))
}

// POST /api/admin/users/:id/reset-password
func (h *AdminHandler) ResetUserPassword(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid user ID"))
		return
	}

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_request", err.Error()))
		return
	}

	err = h.adminService.ResetUserPassword(c.Request.Context(), userID, adminID, req.NewPassword)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Password reset successfully",
	}))
}

// GET /api/admin/export/:dataType
func (h *AdminHandler) ExportData(c *gin.Context) {
	dataType := c.Param("dataType")
	format := c.Query("format")
	if format == "" {
		format = "csv"
	}

	// Validate data type
	validDataTypes := map[string]bool{
		"users":         true,
		"communities":   true,
		"reports":       true,
		"announcements": true,
		"events":        true,
	}

	if !validDataTypes[dataType] {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_data_type", "Invalid data type for export"))
		return
	}

	// Set appropriate headers for file download
	timestamp := time.Now().Format("2006-01-02-150405")
	filename := dataType + "-export-" + timestamp + "." + format
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename="+filename)

	// For now, return a basic CSV header
	// TODO: Implement actual data export from service layer
	csv := "id,name,created_at\n"
	c.String(http.StatusOK, csv)
}

// Legacy handler for compatibility
func AdminPanel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "admin ok"})
}
