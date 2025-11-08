package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/connect-univyn/connect-server/internal/service/admin"
	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	adminService *admin.Service
}

func NewAdminHandler(adminService *admin.Service) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}


func (h *AdminHandler) SuspendUser(c *gin.Context) {
	
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "No auth payload"))
		return
	}
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	
	var req struct {
		Reason       string `json:"reason" binding:"required"`
		Notes        string `json:"notes"`
		DurationDays int    `json:"duration_days"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	
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


func (h *AdminHandler) UnsuspendUser(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
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


func (h *AdminHandler) BanUser(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) GetReports(c *gin.Context) {
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
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


func (h *AdminHandler) GetSpaceActivities(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
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


func (h *AdminHandler) GetDashboardStats(c *gin.Context) {
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
		return
	}

	stats, err := h.adminService.GetDashboardStats(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(stats))
}


func (h *AdminHandler) GetUsers(c *gin.Context) {
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
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


func (h *AdminHandler) DeleteUser(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
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


func (h *AdminHandler) GetTutorApplications(c *gin.Context) {
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
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


func (h *AdminHandler) ApproveTutorApplication(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid application ID"))
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) RejectTutorApplication(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid application ID"))
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) GetMentorApplications(c *gin.Context) {
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
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


func (h *AdminHandler) ApproveMentorApplication(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid application ID"))
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) RejectMentorApplication(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid application ID"))
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) ResolveReport(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	reportID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid report ID"))
		return
	}

	var req struct {
		Action string `json:"action"`
		Notes  string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) EscalateReport(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	reportID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid report ID"))
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


func (h *AdminHandler) GetGroups(c *gin.Context) {
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
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


func (h *AdminHandler) ApproveGroup(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid group ID"))
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


func (h *AdminHandler) RejectGroup(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid group ID"))
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) DeleteGroup(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid group ID"))
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
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	
	valueJSON, err := json.Marshal(req.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid value format"))
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




func (h *AdminHandler) GetUserGrowth(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
		return
	}

	
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


func (h *AdminHandler) GetEngagementMetrics(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
		return
	}

	since := time.Now().AddDate(0, -1, 0) 
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


func (h *AdminHandler) GetActivityAnalytics(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
		return
	}

	since := time.Now().AddDate(0, -1, 0) 
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


func (h *AdminHandler) UpdateAdminRole(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	var req struct {
		Roles []string `json:"roles" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) UpdateAdminStatus(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminUserID, _ := uuid.Parse(authPayload.UserID)

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) MarkNotificationRead(c *gin.Context) {
	notificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid notification ID"))
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


func (h *AdminHandler) DeleteNotification(c *gin.Context) {
	notificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid notification ID"))
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




func (h *AdminHandler) GetCommunities(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
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


func (h *AdminHandler) CreateCommunity(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	var req admin.CreateCommunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) UpdateCommunity(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid community ID"))
		return
	}

	var req admin.UpdateCommunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) DeleteCommunity(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid community ID"))
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


func (h *AdminHandler) UpdateCommunityStatus(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid community ID"))
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) AssignCommunityModerator(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid community ID"))
		return
	}

	var req struct {
		UserID      string   `json:"user_id" binding:"required"`
		Permissions []string `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
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




func (h *AdminHandler) GetAnnouncements(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
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


func (h *AdminHandler) CreateAnnouncement(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	var req admin.CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) UpdateAnnouncement(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	announcementID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid announcement ID"))
		return
	}

	var req admin.UpdateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) DeleteAnnouncement(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	announcementID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid announcement ID"))
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


func (h *AdminHandler) UpdateAnnouncementStatus(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	announcementID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid announcement ID"))
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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




func (h *AdminHandler) GetEvents(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
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


func (h *AdminHandler) CreateEvent(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	var req admin.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) UpdateEvent(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid event ID"))
		return
	}

	var req admin.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) DeleteEvent(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid event ID"))
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


func (h *AdminHandler) UpdateEventStatus(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid event ID"))
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) GetEventRegistrations(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid event ID"))
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




func (h *AdminHandler) CreateUser(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	var req admin.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) UpdateUser(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	var req admin.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) ResetUserPassword(c *gin.Context) {
	payload, _ := c.Get("authorization_payload")
	authPayload := payload.(*auth.Payload)
	adminID, _ := uuid.Parse(authPayload.UserID)

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
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


func (h *AdminHandler) ExportData(c *gin.Context) {
	dataType := c.Param("dataType")
	format := c.Query("format")
	if format == "" {
		format = "csv"
	}

	
	validDataTypes := map[string]bool{
		"users":         true,
		"communities":   true,
		"reports":       true,
		"announcements": true,
		"events":        true,
	}

	if !validDataTypes[dataType] {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid data type for export"))
		return
	}

	
	timestamp := time.Now().Format("2006-01-02-150405")
	filename := dataType + "-export-" + timestamp + "." + format
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename="+filename)

	
	
	csv := "id,name,created_at\n"
	c.String(http.StatusOK, csv)
}


func AdminPanel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "admin ok"})
}
