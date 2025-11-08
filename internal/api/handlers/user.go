package handlers

import (
	"net/http"
	"strconv"

	"github.com/connect-univyn/connect_server/internal/service/users"
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *users.Service
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *users.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser handles POST /api/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req users.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(user))
}

// GetUser handles GET /api/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid user ID format"))
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(user))
}

// GetUserByUsername handles GET /api/users/username/:username
func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	spaceIDStr := c.Query("space_id")

	if spaceIDStr == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("missing_space_id", "space_id query parameter is required"))
		return
	}

	user, err := h.userService.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(user))
}

// UpdateUser handles PUT /api/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid user ID format"))
		return
	}

	// TODO: Check if the authenticated user has permission to update this profile
	// authPayload := c.MustGet("authorization_payload").(*auth.Payload)

	var req users.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), userID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(user))
}

// UpdatePassword handles PUT /api/users/:id/password
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid user ID format"))
		return
	}

	// TODO: Check if the authenticated user is updating their own password
	// authPayload := c.MustGet("authorization_payload").(*auth.Payload)

	var req users.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	err = h.userService.UpdatePassword(c.Request.Context(), userID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Password updated successfully"}))
}

// DeactivateUser handles DELETE /api/users/:id
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid user ID format"))
		return
	}

	// TODO: Check if the authenticated user has permission to deactivate this account
	// authPayload := c.MustGet("authorization_payload").(*auth.Payload)

	err = h.userService.DeactivateUser(c.Request.Context(), userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "User deactivated successfully"}))
}

// SearchUsers handles GET /api/users/search
func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	spaceIDStr := c.Query("space_id")

	if query == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("missing_query", "Search query (q) is required"))
		return
	}

	if spaceIDStr == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("missing_space_id", "space_id query parameter is required"))
		return
	}

	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID format"))
		return
	}

	users, err := h.userService.SearchUsers(c.Request.Context(), query, spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(users))
}

// GetSuggestedUsers handles GET /api/users/suggested
func (h *UserHandler) GetSuggestedUsers(c *gin.Context) {
	// Get authenticated user ID from context
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Authentication required"))
		return
	}
	authPayload := payload.(*auth.Payload)

	parsedUserID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user_id", "Invalid user ID"))
		return
	}

	// Get space_id from query parameters
	spaceIDStr := c.Query("space_id")
	if spaceIDStr == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("missing_space_id", "space_id query parameter is required"))
		return
	}

	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID format"))
		return
	}

	// Get pagination parameters
	page := int32(1)
	if pageStr := c.Query("page"); pageStr != "" {
		pageInt, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil || pageInt < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid page number"))
			return
		}
		page = int32(pageInt)
	}

	limit := int32(10)
	if limitStr := c.Query("limit"); limitStr != "" {
		limitInt, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || limitInt < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid limit"))
			return
		}
		limit = int32(limitInt)
	}

	// Calculate offset
	offset := (page - 1) * limit

	users, err := h.userService.GetSuggestedUsers(c.Request.Context(), parsedUserID, spaceID, limit, offset)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(users))
}

// FollowUser handles POST /api/users/:id/follow
func (h *UserHandler) FollowUser(c *gin.Context) {
	// Get authenticated user ID
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Authentication required"))
		return
	}
	authPayload := payload.(*auth.Payload)

	followerID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user_id", "Invalid user ID"))
		return
	}

	// Get user to follow from URL parameter
	followingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user_id", "Invalid user ID format"))
		return
	}

	// Get space_id from query parameters
	spaceIDStr := c.Query("space_id")
	if spaceIDStr == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("missing_space_id", "space_id query parameter is required"))
		return
	}

	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_space_id", "Invalid space ID format"))
		return
	}

	// Follow user
	if err := h.userService.FollowUser(c.Request.Context(), followerID, followingID, spaceID); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Successfully followed user"}))
}

// UnfollowUser handles DELETE /api/users/:id/follow
func (h *UserHandler) UnfollowUser(c *gin.Context) {
	// Get authenticated user ID
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Authentication required"))
		return
	}
	authPayload := payload.(*auth.Payload)

	followerID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user_id", "Invalid user ID"))
		return
	}

	// Get user to unfollow from URL parameter
	followingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user_id", "Invalid user ID format"))
		return
	}

	// Unfollow user
	if err := h.userService.UnfollowUser(c.Request.Context(), followerID, followingID); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Successfully unfollowed user"}))
}

// CheckIfFollowing handles GET /api/users/:id/following/status
func (h *UserHandler) CheckIfFollowing(c *gin.Context) {
	// Get authenticated user ID
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Authentication required"))
		return
	}
	authPayload := payload.(*auth.Payload)

	followerID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user_id", "Invalid user ID"))
		return
	}

	// Get user to check from URL parameter
	followingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user_id", "Invalid user ID format"))
		return
	}

	// Check if following
	isFollowing, err := h.userService.CheckIfFollowing(c.Request.Context(), followerID, followingID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"is_following": isFollowing}))
}

// GetFollowers handles GET /api/users/:id/followers
func (h *UserHandler) GetFollowers(c *gin.Context) {
	// Get user ID from URL parameter
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user_id", "Invalid user ID format"))
		return
	}

	// Get pagination parameters
	page := int32(1)
	if pageStr := c.Query("page"); pageStr != "" {
		pageInt, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil || pageInt < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid page number"))
			return
		}
		page = int32(pageInt)
	}

	limit := int32(20)
	if limitStr := c.Query("limit"); limitStr != "" {
		limitInt, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || limitInt < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid limit"))
			return
		}
		limit = int32(limitInt)
	}

	// Get followers
	followers, err := h.userService.GetFollowers(c.Request.Context(), userID, page, limit)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(followers))
}

// GetFollowing handles GET /api/users/:id/following
func (h *UserHandler) GetFollowing(c *gin.Context) {
	// Get user ID from URL parameter
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user_id", "Invalid user ID format"))
		return
	}

	// Get pagination parameters
	page := int32(1)
	if pageStr := c.Query("page"); pageStr != "" {
		pageInt, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil || pageInt < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid page number"))
			return
		}
		page = int32(pageInt)
	}

	limit := int32(20)
	if limitStr := c.Query("limit"); limitStr != "" {
		limitInt, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || limitInt < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid limit"))
			return
		}
		limit = int32(limitInt)
	}

	// Get following
	following, err := h.userService.GetFollowing(c.Request.Context(), userID, page, limit)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(following))
}
