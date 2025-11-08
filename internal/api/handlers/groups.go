package handlers

import (
	"net/http"

	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/connect-univyn/connect-server/internal/service/groups"
	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type GroupHandler struct {
	groupService *groups.Service
}


func NewGroupHandler(groupService *groups.Service) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}


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


func (h *GroupHandler) GetGroup(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID format"))
		return
	}
	
	
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


func (h *GroupHandler) ListGroups(c *gin.Context) {
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID format"))
		return
	}
	
	
	var userID uuid.UUID
	if payload, exists := c.Get("authorization_payload"); exists {
		authPayload := payload.(*auth.Payload)
		userID, _ = uuid.Parse(authPayload.UserID)
	}
	
	page, limit := parsePagination(c)
	sortBy := c.DefaultQuery("sort", "recent") 
	
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
