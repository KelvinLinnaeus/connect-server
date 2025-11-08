package websocket

import (
	"net/http"

	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin checking based on CORS configuration
		// For now, allow all origins (to be secured in production)
		return true
	},
}

// Handler handles WebSocket upgrade requests
type Handler struct {
	manager    *Manager
	tokenMaker auth.Maker
}

// NewHandler creates a new WebSocket handler
func NewHandler(manager *Manager, tokenMaker auth.Maker) *Handler {
	return &Handler{
		manager:    manager,
		tokenMaker: tokenMaker,
	}
}

// HandleWebSocket handles WebSocket connection requests
func (h *Handler) HandleWebSocket(c *gin.Context) {
	// Get token from query parameter or header
	token := c.Query("token")
	if token == "" {
		token = c.GetHeader("Authorization")
		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Authentication token required"))
		return
	}

	// Verify token
	payload, err := h.tokenMaker.VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("invalid_token", "Invalid or expired token"))
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("invalid_token", "Invalid user ID in token"))
		return
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade WebSocket connection")
		return
	}

	// Get client IP address (handles X-Forwarded-For properly)
	ipAddress := c.ClientIP()

	// Create client
	client := NewClient(conn, userID, ipAddress, h.manager)

	// Add space ID from context if available
	if spaceIDStr := c.Query("space_id"); spaceIDStr != "" {
		if spaceID, err := uuid.Parse(spaceIDStr); err == nil {
			client.SpaceID = &spaceID
		}
	}

	// Register client
	h.manager.Register(client)

	log.Info().
		Str("client_id", client.ID).
		Str("user_id", userID.String()).
		Str("remote_addr", c.Request.RemoteAddr).
		Msg("WebSocket connection established")

	// Start client pumps
	go client.WritePump()
	go client.ReadPump()

	// Auto-subscribe to user's personal channel
	client.SubscriptionsMu.Lock()
	userChannel := Channel.User(userID)
	client.Subscriptions[userChannel] = true
	client.SubscriptionsMu.Unlock()

	// Send welcome message
	client.sendMessage(ServerMessage{
		Type:    MessageTypeAck,
		Channel: "",
		Payload: map[string]interface{}{
			"message":       "Connected successfully",
			"client_id":     client.ID,
			"subscriptions": []string{userChannel},
		},
	})
}

// HandleMetrics handles metrics endpoint
func (h *Handler) HandleMetrics(c *gin.Context) {
	metrics := h.manager.GetMetrics()

	c.JSON(http.StatusOK, gin.H{
		"active_connections":     metrics.ActiveConnections,
		"total_connections":      metrics.TotalConnections,
		"connections_rejected":   metrics.ConnectionsRejected,
		"messages_received":      metrics.MessagesReceived,
		"messages_sent":          metrics.MessagesSent,
		"errors":                 metrics.Errors,
		"last_error":             metrics.LastError,
		"last_error_time":        metrics.LastErrorTime,
		"average_latency_ms":     metrics.GetAverageLatencyMs(),
		"message_throughput_sec": metrics.GetMessageThroughput(),
	})
}

// HandlePresence handles presence check endpoint
func (h *Handler) HandlePresence(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_request", "User ID required"))
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_request", "Invalid user ID"))
		return
	}

	isOnline := h.manager.IsUserOnline(userID)
	connections := h.manager.GetUserConnections(userID)

	c.JSON(http.StatusOK, gin.H{
		"user_id":           userID,
		"online":            isOnline,
		"connection_count":  len(connections),
	})
}

// HandleBulkPresence handles bulk presence check endpoint
func (h *Handler) HandleBulkPresence(c *gin.Context) {
	var req struct {
		UserIDs []string `json:"user_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_request", err.Error()))
		return
	}

	presence := make(map[string]bool)
	for _, userIDStr := range req.UserIDs {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			continue
		}
		presence[userIDStr] = h.manager.IsUserOnline(userID)
	}

	c.JSON(http.StatusOK, gin.H{
		"presence": presence,
	})
}
