package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	"github.com/connect-univyn/connect_server/internal/service/users"
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	userService          *users.Service
	tokenMaker           auth.Maker
	store                db.Store
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	userService *users.Service,
	tokenMaker auth.Maker,
	store db.Store,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
) *AuthHandler {
	return &AuthHandler{
		userService:          userService,
		tokenMaker:           tokenMaker,
		store:                store,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	AccessToken           string              `json:"access_token"`
	RefreshToken          string              `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time           `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time           `json:"refresh_token_expires_at"`
	User                  *users.UserResponse `json:"user"`
}

// Login handles POST /api/users/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	req.Email = strings.ToLower(req.Email)

	// Authenticate user
	user, err := h.userService.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// Create access token
	accessToken, accessPayload, err := h.tokenMaker.CreateToken(
		user.ID.String(),
		user.Username,
		user.SpaceID.String(),
		h.accessTokenDuration,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse("token_error", "Failed to create access token"))
		return
	}

	// Create refresh token
	refreshToken, refreshPayload, err := h.tokenMaker.CreateToken(
		user.ID.String(),
		user.Username,
		user.SpaceID.String(),
		h.refreshTokenDuration,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse("token_error", "Failed to create refresh token"))
		return
	}

	// Create session in database
	var ipAddress sql.NullString
	if c.ClientIP() != "" {
		ipAddress = sql.NullString{String: c.ClientIP(), Valid: true}
	}

	_, err = h.store.CreateSession(c.Request.Context(), db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		IpAddress:    ipAddress,
		IsBlocked:    false,
		SpaceID:      user.SpaceID,
		LastActivity: sql.NullTime{Time: time.Now(), Valid: true},
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse("session_error", "Failed to create session"))
		return
	}

	response := LoginResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  user,
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(response))
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse represents refresh token response
type RefreshTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// RefreshToken handles POST /api/users/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	// Verify refresh token
	refreshPayload, err := h.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// Get session from database
	session, err := h.store.GetSession(c.Request.Context(), refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, util.NewErrorResponse("invalid_session", "Session not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse("session_error", "Failed to get session"))
		return
	}

	// Check if session is blocked
	if session.IsBlocked {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("blocked_session", "Session is blocked"))
		return
	}

	// Check if session user matches token user
	if session.UserID.String() != refreshPayload.UserID {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("invalid_session", "Session user mismatch"))
		return
	}

	// Check if refresh token matches
	if session.RefreshToken != req.RefreshToken {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("invalid_token", "Token mismatch"))
		return
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("expired_session", "Session has expired"))
		return
	}

	// Create new access token
	accessToken, accessPayload, err := h.tokenMaker.CreateToken(
		refreshPayload.UserID,
		refreshPayload.Username,
		refreshPayload.SpaceID,
		h.accessTokenDuration,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse("token_error", "Failed to create access token"))
		return
	}

	response := RefreshTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(response))
}

// Logout handles POST /api/users/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get authenticated user from context (set by auth middleware)
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}

	authPayload := payload.(*auth.Payload)

	// TODO: Block or delete the session
	// For now, we just return success
	// In production, you'd want to:
	// 1. Mark the session as blocked in the database
	// 2. Or delete the session entirely
	// 3. Consider implementing a token blacklist

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Logged out successfully",
		"user_id": authPayload.UserID,
	}))
}

// GetSession handles GET /api/sessions/:id
func (h *AuthHandler) GetSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid session ID format"))
		return
	}

	// TODO: Verify user has permission to view this session

	session, err := h.store.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.NewErrorResponse("not_found", "Session not found"))
			return
		}
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(session))
}
