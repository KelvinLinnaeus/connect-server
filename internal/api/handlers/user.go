package handlers

import (
	"net/http"
	"strconv"

	"github.com/connect-univyn/connect-server/internal/service/users"
	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type UserHandler struct {
	userService *users.Service
}


func NewUserHandler(userService *users.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}


func (h *UserHandler) CreateUser(c *gin.Context) {
	var req users.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(user))
}


func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID format"))
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(user))
}


func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	spaceIDStr := c.Query("space_id")

	if spaceIDStr == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "space_id query parameter is required"))
		return
	}

	user, err := h.userService.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(user))
}


func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID format"))
		return
	}

	
	

	var req users.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), userID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(user))
}


func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID format"))
		return
	}

	
	

	var req users.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	err = h.userService.UpdatePassword(c.Request.Context(), userID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Password updated successfully"}))
}


func (h *UserHandler) DeactivateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID format"))
		return
	}

	
	

	err = h.userService.DeactivateUser(c.Request.Context(), userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "User deactivated successfully"}))
}


func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	spaceIDStr := c.Query("space_id")

	if query == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Search query (q) is required"))
		return
	}

	if spaceIDStr == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "space_id query parameter is required"))
		return
	}

	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID format"))
		return
	}

	users, err := h.userService.SearchUsers(c.Request.Context(), query, spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(users))
}


func (h *UserHandler) GetSuggestedUsers(c *gin.Context) {
	
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Authentication required"))
		return
	}
	authPayload := payload.(*auth.Payload)

	parsedUserID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	
	spaceIDStr := c.Query("space_id")
	if spaceIDStr == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "space_id query parameter is required"))
		return
	}

	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID format"))
		return
	}

	
	page := int32(1)
	if pageStr := c.Query("page"); pageStr != "" {
		pageInt, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil || pageInt < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid page number"))
			return
		}
		page = int32(pageInt)
	}

	limit := int32(10)
	if limitStr := c.Query("limit"); limitStr != "" {
		limitInt, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || limitInt < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid limit"))
			return
		}
		limit = int32(limitInt)
	}

	
	offset := (page - 1) * limit

	users, err := h.userService.GetSuggestedUsers(c.Request.Context(), parsedUserID, spaceID, limit, offset)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(users))
}


func (h *UserHandler) FollowUser(c *gin.Context) {
	
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Authentication required"))
		return
	}
	authPayload := payload.(*auth.Payload)

	followerID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	
	followingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID format"))
		return
	}

	
	spaceIDStr := c.Query("space_id")
	if spaceIDStr == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "space_id query parameter is required"))
		return
	}

	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID format"))
		return
	}

	
	if err := h.userService.FollowUser(c.Request.Context(), followerID, followingID, spaceID); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Successfully followed user"}))
}


func (h *UserHandler) UnfollowUser(c *gin.Context) {
	
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Authentication required"))
		return
	}
	authPayload := payload.(*auth.Payload)

	followerID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	
	followingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID format"))
		return
	}

	
	if err := h.userService.UnfollowUser(c.Request.Context(), followerID, followingID); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Successfully unfollowed user"}))
}


func (h *UserHandler) CheckIfFollowing(c *gin.Context) {
	
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Authentication required"))
		return
	}
	authPayload := payload.(*auth.Payload)

	followerID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	
	followingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID format"))
		return
	}

	
	isFollowing, err := h.userService.CheckIfFollowing(c.Request.Context(), followerID, followingID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"is_following": isFollowing}))
}


func (h *UserHandler) GetFollowers(c *gin.Context) {
	
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID format"))
		return
	}

	
	page := int32(1)
	if pageStr := c.Query("page"); pageStr != "" {
		pageInt, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil || pageInt < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid page number"))
			return
		}
		page = int32(pageInt)
	}

	limit := int32(20)
	if limitStr := c.Query("limit"); limitStr != "" {
		limitInt, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || limitInt < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid limit"))
			return
		}
		limit = int32(limitInt)
	}

	
	followers, err := h.userService.GetFollowers(c.Request.Context(), userID, page, limit)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(followers))
}


func (h *UserHandler) GetFollowing(c *gin.Context) {
	
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID format"))
		return
	}

	
	page := int32(1)
	if pageStr := c.Query("page"); pageStr != "" {
		pageInt, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil || pageInt < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid page number"))
			return
		}
		page = int32(pageInt)
	}

	limit := int32(20)
	if limitStr := c.Query("limit"); limitStr != "" {
		limitInt, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || limitInt < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid limit"))
			return
		}
		limit = int32(limitInt)
	}

	
	following, err := h.userService.GetFollowing(c.Request.Context(), userID, page, limit)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(following))
}
