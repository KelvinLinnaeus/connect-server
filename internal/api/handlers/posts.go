package handlers

import (
	"net/http"
	"strconv"

	"github.com/connect-univyn/connect-server/internal/service/posts"
	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type PostHandler struct {
	postService *posts.Service
}


func NewPostHandler(postService *posts.Service) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}


func (h *PostHandler) CreatePost(c *gin.Context) {
	
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	var req posts.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	
	authorID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_user", "Invalid user ID"))
		return
	}
	req.AuthorID = authorID

	post, err := h.postService.CreatePost(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(post))
}


func (h *PostHandler) GetPost(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid post ID format"))
		return
	}

	
	var userID uuid.UUID
	if payload, exists := c.Get("authorization_payload"); exists {
		authPayload := payload.(*auth.Payload)
		userID, _ = uuid.Parse(authPayload.UserID)
	}

	post, err := h.postService.GetPostByID(c.Request.Context(), postID, userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	
	go h.postService.IncrementPostViews(c.Request.Context(), postID)

	c.JSON(http.StatusOK, util.NewSuccessResponse(post))
}


func (h *PostHandler) DeletePost(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid post ID format"))
		return
	}

	
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	authorID, _ := uuid.Parse(authPayload.UserID)

	err = h.postService.DeletePost(c.Request.Context(), postID, authorID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Post deleted successfully"}))
}


func (h *PostHandler) GetUserPosts(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid user ID format"))
		return
	}

	
	limit, offset := parsePagination(c)

	posts, err := h.postService.GetUserPosts(c.Request.Context(), userID, int32(limit), int32(offset))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(posts))
}


func (h *PostHandler) GetUserFeed(c *gin.Context) {
	
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	userID, _ := uuid.Parse(authPayload.UserID)
	spaceID, _ := uuid.Parse(authPayload.SpaceID)

	
	limit, offset := parsePagination(c)

	posts, err := h.postService.GetUserFeed(c.Request.Context(), userID, spaceID, int32(limit), int32(offset))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(posts))
}


func (h *PostHandler) GetCommunityPosts(c *gin.Context) {
	communityID, err := uuid.Parse(c.Param("community_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid community ID format"))
		return
	}

	
	var userID uuid.UUID
	if payload, exists := c.Get("authorization_payload"); exists {
		authPayload := payload.(*auth.Payload)
		userID, _ = uuid.Parse(authPayload.UserID)
	}

	
	limit, offset := parsePagination(c)

	posts, err := h.postService.GetCommunityPosts(c.Request.Context(), userID, communityID, int32(limit), int32(offset))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(posts))
}


func (h *PostHandler) GetGroupPosts(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid group ID format"))
		return
	}

	
	var userID uuid.UUID
	if payload, exists := c.Get("authorization_payload"); exists {
		authPayload := payload.(*auth.Payload)
		userID, _ = uuid.Parse(authPayload.UserID)
	}

	
	limit, offset := parsePagination(c)

	posts, err := h.postService.GetGroupPosts(c.Request.Context(), userID, groupID, int32(limit), int32(offset))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(posts))
}


func (h *PostHandler) GetTrendingPosts(c *gin.Context) {
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

	posts, err := h.postService.GetTrendingPosts(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(posts))
}


func (h *PostHandler) GetTrendingTopics(c *gin.Context) {
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

	
	offset := (page - 1) * limit

	topics, err := h.postService.GetTrendingTopics(c.Request.Context(), spaceID, limit, offset)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(topics))
}


func (h *PostHandler) SearchPosts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("missing_query", "Search query (q) is required"))
		return
	}

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

	
	limit, offset := parsePagination(c)

	posts, err := h.postService.SearchPosts(c.Request.Context(), query, spaceID, int32(limit), int32(offset))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(posts))
}


func (h *PostHandler) AdvancedSearchPosts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("missing_query", "Search query (q) is required"))
		return
	}

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

	
	limit, offset := parsePagination(c)

	posts, err := h.postService.AdvancedSearchPosts(c.Request.Context(), query, spaceID, int32(limit), int32(offset))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(posts))
}


func (h *PostHandler) GetUserLikedPosts(c *gin.Context) {
	
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	userID, _ := uuid.Parse(authPayload.UserID)

	
	limit, offset := parsePagination(c)

	posts, err := h.postService.GetUserLikedPosts(c.Request.Context(), userID, int32(limit), int32(offset))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(posts))
}


func (h *PostHandler) GetPostComments(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid post ID format"))
		return
	}

	comments, err := h.postService.GetPostComments(c.Request.Context(), postID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(comments))
}


func (h *PostHandler) GetPostLikes(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid post ID format"))
		return
	}

	likes, err := h.postService.GetPostLikes(c.Request.Context(), postID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(likes))
}


func (h *PostHandler) CreateComment(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid post ID format"))
		return
	}

	
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	authorID, _ := uuid.Parse(authPayload.UserID)

	var req posts.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	req.PostID = postID
	req.AuthorID = authorID

	comment, err := h.postService.CreateComment(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(comment))
}


func (h *PostHandler) CreateRepost(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid post ID format"))
		return
	}

	
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	authorID, _ := uuid.Parse(authPayload.UserID)
	spaceID, _ := uuid.Parse(authPayload.SpaceID)

	var req posts.CreateRepostParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	req.AuthorID = authorID
	req.SpaceID = spaceID
	req.QuotedPostID = &postID

	repost, err := h.postService.CreateRepost(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(repost))
}


func (h *PostHandler) TogglePostLike(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid post ID format"))
		return
	}

	
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	userID, _ := uuid.Parse(authPayload.UserID)

	likesCount, err := h.postService.TogglePostLike(c.Request.Context(), userID, postID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"likes_count": likesCount,
	}))
}


func (h *PostHandler) ToggleCommentLike(c *gin.Context) {
	commentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid comment ID format"))
		return
	}

	
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	userID, _ := uuid.Parse(authPayload.UserID)

	liked, err := h.postService.ToggleCommentLike(c.Request.Context(), userID, commentID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"liked": liked,
	}))
}


func (h *PostHandler) PinPost(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid post ID format"))
		return
	}

	var req struct {
		IsPinned bool `json:"is_pinned"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	

	err = h.postService.PinPost(c.Request.Context(), postID, req.IsPinned)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message":   "Post pin status updated",
		"is_pinned": req.IsPinned,
	}))
}


func parsePagination(c *gin.Context) (int, int) {
	limit := 20 
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			offset = (page - 1) * limit
		}
	}

	return limit, offset
}
