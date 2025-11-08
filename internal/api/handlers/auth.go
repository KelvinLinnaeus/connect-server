package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	db "github.com/connect-univyn/connect-server/db/sqlc"
	"github.com/connect-univyn/connect-server/internal/service/users"
	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type AuthHandler struct {
	userService          *users.Service
	tokenMaker           auth.Maker
	store                db.Store
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}


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


type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}


type LoginResponse struct {
	AccessToken           string              `json:"access_token"`
	RefreshToken          string              `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time           `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time           `json:"refresh_token_expires_at"`
	User                  *users.UserResponse `json:"user"`
}


func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	req.Email = strings.ToLower(req.Email)

	
	user, err := h.userService.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	
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


type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}


type RefreshTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}


func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	
	refreshPayload, err := h.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	
	session, err := h.store.GetSession(c.Request.Context(), refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, util.NewErrorResponse("invalid_session", "Session not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, util.NewErrorResponse("session_error", "Failed to get session"))
		return
	}

	
	if session.IsBlocked {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("blocked_session", "Session is blocked"))
		return
	}

	
	if session.UserID.String() != refreshPayload.UserID {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("invalid_session", "Session user mismatch"))
		return
	}

	
	if session.RefreshToken != req.RefreshToken {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("invalid_token", "Token mismatch"))
		return
	}

	
	if time.Now().After(session.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("expired_session", "Session has expired"))
		return
	}

	
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


func (h *AuthHandler) Logout(c *gin.Context) {
	
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}

	authPayload := payload.(*auth.Payload)

	
	
	
	
	
	

	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{
		"message": "Logged out successfully",
		"user_id": authPayload.UserID,
	}))
}


func (h *AuthHandler) GetSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid session ID format"))
		return
	}

	

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
