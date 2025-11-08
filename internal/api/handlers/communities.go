package handlers

import (
	"net/http"

	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/connect-univyn/connect_server/internal/service/communities"
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CommunityHandler handles community-related HTTP requests
type CommunityHandler struct {
	communityService *communities.Service
}

// NewCommunityHandler creates a new community handler
func NewCommunityHandler(communityService *communities.Service) *CommunityHandler {
	return &CommunityHandler{
		communityService: communityService,
	}
}

// CreateCommunity handles POST /api/communities
func (h *CommunityHandler) CreateCommunity(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	var req communities.CreateCommunityRequest
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
	
	community, err := h.communityService.CreateCommunity(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, util.NewSuccessResponse(community))
}

// GetCommunity handles GET /api/communities/:id
func (h *CommunityHandler) GetCommunity(c *gin.Context) {
	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid community ID format"))
		return
	}
	
	// Try to get user ID from auth payload (optional for public communities)
	var userID uuid.UUID
	if payload, exists := c.Get("authorization_payload"); exists {
		authPayload := payload.(*auth.Payload)
		userID, _ = uuid.Parse(authPayload.UserID)
	}
	
	community, err := h.communityService.GetCommunityByID(c.Request.Context(), userID, communityID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(community))
}

// GetCommunityBySlug handles GET /api/communities/slug/:slug
func (h *CommunityHandler) GetCommunityBySlug(c *gin.Context) {
	slug := c.Param("slug")
	spaceIDStr := c.Query("space_id")
	
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID format"))
		return
	}
	
	// Try to get user ID from auth payload (optional for public communities)
	var userID uuid.UUID
	if payload, exists := c.Get("authorization_payload"); exists {
		authPayload := payload.(*auth.Payload)
		userID, _ = uuid.Parse(authPayload.UserID)
	}
	
	community, err := h.communityService.GetCommunityBySlug(c.Request.Context(), userID, spaceID, slug)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(community))
}

// ListCommunities handles GET /api/communities
func (h *CommunityHandler) ListCommunities(c *gin.Context) {
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
		spaceID, _ = uuid.Parse(authPayload.SpaceID)
	}
	
	page, limit := parsePagination(c)
	sortBy := c.DefaultQuery("sort", "recent") // members, posts, recent
	
	params := communities.ListCommunitiesParams{
		UserID:  userID,
		SpaceID: spaceID,
		SortBy:  sortBy,
		Page:    int32(page),
		Limit:   int32(limit),
	}
	
	communityList, err := h.communityService.ListCommunities(c.Request.Context(), params)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(communityList))
}

// SearchCommunities handles GET /api/communities/search
func (h *CommunityHandler) SearchCommunities(c *gin.Context) {
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
	
	params := communities.SearchCommunitiesParams{
		UserID:  userID,
		SpaceID: spaceID,
		Query:   query,
	}
	
	results, err := h.communityService.SearchCommunities(c.Request.Context(), params)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(results))
}

// UpdateCommunity handles PUT /api/communities/:id
func (h *CommunityHandler) UpdateCommunity(c *gin.Context) {
	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid community ID format"))
		return
	}
	
	var req communities.UpdateCommunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	
	community, err := h.communityService.UpdateCommunity(c.Request.Context(), communityID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(community))
}

// GetCommunityMembers handles GET /api/communities/:id/members
func (h *CommunityHandler) GetCommunityMembers(c *gin.Context) {
	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid community ID format"))
		return
	}
	
	members, err := h.communityService.GetCommunityMembers(c.Request.Context(), communityID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(members))
}

// GetCommunityModerators handles GET /api/communities/:id/moderators
func (h *CommunityHandler) GetCommunityModerators(c *gin.Context) {
	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid community ID format"))
		return
	}
	
	moderators, err := h.communityService.GetCommunityModerators(c.Request.Context(), communityID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(moderators))
}

// GetCommunityAdmins handles GET /api/communities/:id/admins
func (h *CommunityHandler) GetCommunityAdmins(c *gin.Context) {
	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid community ID format"))
		return
	}
	
	admins, err := h.communityService.GetCommunityAdmins(c.Request.Context(), communityID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(admins))
}

// JoinCommunity handles POST /api/communities/:id/join
func (h *CommunityHandler) JoinCommunity(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid community ID format"))
		return
	}
	
	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user", "Invalid user ID"))
		return
	}
	
	membership, err := h.communityService.JoinCommunity(c.Request.Context(), communityID, userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(membership))
}

// LeaveCommunity handles POST /api/communities/:id/leave
func (h *CommunityHandler) LeaveCommunity(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid community ID format"))
		return
	}
	
	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user", "Invalid user ID"))
		return
	}
	
	err = h.communityService.LeaveCommunity(c.Request.Context(), communityID, userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Successfully left community"}))
}

// AddCommunityModerator handles POST /api/communities/:id/moderators
func (h *CommunityHandler) AddCommunityModerator(c *gin.Context) {
	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid community ID format"))
		return
	}
	
	var req communities.AddModeratorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	
	membership, err := h.communityService.AddCommunityModerator(c.Request.Context(), communityID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(membership))
}

// RemoveCommunityModerator handles DELETE /api/communities/:id/moderators/:userId
func (h *CommunityHandler) RemoveCommunityModerator(c *gin.Context) {
	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid community ID format"))
		return
	}
	
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user_id", "Invalid user ID format"))
		return
	}
	
	err = h.communityService.RemoveCommunityModerator(c.Request.Context(), communityID, userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Moderator removed successfully"}))
}

// GetUserCommunities handles GET /api/users/communities
func (h *CommunityHandler) GetUserCommunities(c *gin.Context) {
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
	
	userCommunities, err := h.communityService.GetUserCommunities(c.Request.Context(), userID, spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(userCommunities))
}

// GetCommunityCategories handles GET /api/communities/categories
func (h *CommunityHandler) GetCommunityCategories(c *gin.Context) {
	spaceIDStr := c.Query("space_id")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID format"))
		return
	}
	
	categories, err := h.communityService.GetCommunityCategories(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(categories))
}
