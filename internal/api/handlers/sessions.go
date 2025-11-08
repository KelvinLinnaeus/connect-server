package handlers

import (
	"net/http"

	"github.com/connect-univyn/connect-server/internal/service/sessions"
	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type SessionHandler struct {
	sessionService *sessions.Service
}


func NewSessionHandler(sessionService *sessions.Service) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
	}
}


func (h *SessionHandler) GetSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid session ID format"))
		return
	}

	session, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(session))
}
