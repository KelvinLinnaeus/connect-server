package handlers

import (
	"net/http"

	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/connect-univyn/connect_server/internal/service/groups"
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GroupHandler handles group-related HTTP requests
type GroupHandler struct {
	groupService *groups.Service
}

// NewGroupHandler creates a new group handler
func NewGroupHandler(groupService *groups.Service) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

// CreateGroup handles POST /api/groups
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	var req groups.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	
	// Set the creator from the auth payload
	creatorID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user", "Invalid user ID"))
		return
	}
	req.CreatedBy = creatorID
	
	group, err := h.groupService.CreateGroup(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, util.NewSuccessResponse(group))
}

// GetGroup handles GET /api/groups/:id
func (h *GroupHandler) GetGroup(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID format"))
		return
	}
	
	// Try to get user ID from auth payload (optional for public groups)
	var userID uuid.UUID
	if payload, exists := c.Get("authorization_payload"); exists {
		authPayload := payload.(*auth.Payload)
		userID, _ = uuid.Parse(authPayload.UserID)
	}
	
	group, err := h.groupService.GetGroupByID(c.Request.Context(), userID, groupID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(group))
}

// ListGroups handles GET /api/groups
func (h *GroupHandler) ListGroups(c *gin.Context) {
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID format"))
		return
	}
	
	// Try to get user ID from auth payload (optional)
	var userID uuid.UUID
	if payload, exists := c.Get("authorization_payload"); exists {
		authPayload := payload.(*auth.Payload)
		userID, _ = uuid.Parse(authPayload.UserID)
	}
	
	page, limit := parsePagination(c)
	sortBy := c.DefaultQuery("sort", "recent") // members, recent
	
	params := groups.ListGroupsParams{
		UserID:  userID,
		SpaceID: spaceID,
		SortBy:  sortBy,
		Page:    int32(page),
		Limit:   int32(limit),
	}
	
	groupList, err := h.groupService.ListGroups(c.Request.Context(), params)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(groupList))
}

// SearchGroups handles GET /api/groups/search
func (h *GroupHandler) SearchGroups(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("missing_query", "Search query is required"))
		return
	}
	
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID format"))
		return
	}
	
	// Try to get user ID from auth payload (optional)
	var userID uuid.UUID
	if payload, exists := c.Get("authorization_payload"); exists {
		authPayload := payload.(*auth.Payload)
		userID, _ = uuid.Parse(authPayload.UserID)
	}
	
	params := groups.SearchGroupsParams{
		UserID:  userID,
		SpaceID: spaceID,
		Query:   query,
	}
	
	results, err := h.groupService.SearchGroups(c.Request.Context(), params)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(results))
}

// UpdateGroup handles PUT /api/groups/:id
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID format"))
		return
	}
	
	var req groups.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	
	group, err := h.groupService.UpdateGroup(c.Request.Context(), groupID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(group))
}

// GetUserGroups handles GET /api/users/groups
func (h *GroupHandler) GetUserGroups(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user", "Invalid user ID"))
		return
	}
	
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID format"))
		return
	}
	
	userGroups, err := h.groupService.GetUserGroups(c.Request.Context(), userID, spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(userGroups))
}

// JoinGroup handles POST /api/groups/:id/join
func (h *GroupHandler) JoinGroup(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID format"))
		return
	}
	
	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user", "Invalid user ID"))
		return
	}
	
	membership, err := h.groupService.JoinGroup(c.Request.Context(), groupID, userID, nil)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(membership))
}

// LeaveGroup handles POST /api/groups/:id/leave
func (h *GroupHandler) LeaveGroup(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID format"))
		return
	}
	
	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user", "Invalid user ID"))
		return
	}
	
	err = h.groupService.LeaveGroup(c.Request.Context(), groupID, userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Successfully left group"}))
}

// GetGroupJoinRequests handles GET /api/groups/:id/join-requests
func (h *GroupHandler) GetGroupJoinRequests(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID format"))
		return
	}
	
	requests, err := h.groupService.GetGroupJoinRequests(c.Request.Context(), groupID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(requests))
}

// AddGroupAdmin handles POST /api/groups/:id/admins
func (h *GroupHandler) AddGroupAdmin(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID format"))
		return
	}
	
	var req groups.AddGroupAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	
	membership, err := h.groupService.AddGroupAdmin(c.Request.Context(), groupID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(membership))
}

// RemoveGroupAdmin handles DELETE /api/groups/:id/admins/:userId
func (h *GroupHandler) RemoveGroupAdmin(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID format"))
		return
	}
	
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user_id", "Invalid user ID format"))
		return
	}
	
	err = h.groupService.RemoveGroupAdmin(c.Request.Context(), groupID, userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Admin removed successfully"}))
}

// AddGroupModerator handles POST /api/groups/:id/moderators
func (h *GroupHandler) AddGroupModerator(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID format"))
		return
	}
	
	var req groups.AddGroupModeratorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	
	membership, err := h.groupService.AddGroupModerator(c.Request.Context(), groupID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(membership))
}

// RemoveGroupModerator handles DELETE /api/groups/:id/moderators/:userId
func (h *GroupHandler) RemoveGroupModerator(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID format"))
		return
	}
	
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user_id", "Invalid user ID format"))
		return
	}
	
	err = h.groupService.RemoveGroupModerator(c.Request.Context(), groupID, userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Moderator removed successfully"}))
}

// UpdateGroupMemberRole handles PUT /api/groups/:id/members/:userId/role
func (h *GroupHandler) UpdateGroupMemberRole(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID format"))
		return
	}
	
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user_id", "Invalid user ID format"))
		return
	}
	
	var req groups.UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	
	err = h.groupService.UpdateGroupMemberRole(c.Request.Context(), groupID, userID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Member role updated successfully"}))
}

// CreateProjectRole handles POST /api/groups/:id/roles
func (h *GroupHandler) CreateProjectRole(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID format"))
		return
	}
	
	var req groups.CreateProjectRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	
	role, err := h.groupService.CreateProjectRole(c.Request.Context(), groupID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, util.NewSuccessResponse(role))
}

// GetProjectRoles handles GET /api/groups/:id/roles
func (h *GroupHandler) GetProjectRoles(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID format"))
		return
	}
	
	roles, err := h.groupService.GetProjectRoles(c.Request.Context(), groupID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(roles))
}

// ApplyForProjectRole handles POST /api/roles/:roleId/apply
func (h *GroupHandler) ApplyForProjectRole(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	roleID, err := uuid.Parse(c.Param("roleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_role_id", "Invalid role ID format"))
		return
	}
	
	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user", "Invalid user ID"))
		return
	}
	
	var req groups.ApplyForRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	
	application, err := h.groupService.ApplyForProjectRole(c.Request.Context(), roleID, userID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, util.NewSuccessResponse(application))
}

// GetRoleApplications handles GET /api/groups/:id/applications
func (h *GroupHandler) GetRoleApplications(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID format"))
		return
	}
	
	applications, err := h.groupService.GetRoleApplications(c.Request.Context(), groupID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(applications))
}
