package websocket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// NewClient creates a new WebSocket client
func NewClient(conn *websocket.Conn, userID uuid.UUID, ipAddress string, manager *Manager) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	now := time.Now()

	return &Client{
		ID:            uuid.New().String(),
		UserID:        userID,
		IPAddress:     ipAddress,
		Conn:          conn,
		Send:          make(chan []byte, SendBufferSize),
		Subscriptions: make(map[string]bool),
		Manager:       manager,
		Context:       ctx,
		Cancel:        cancel,
		LastActivity:  now,
		ConnectedAt:   now,
		Metadata:      make(map[string]string),
	}
}

// ReadPump pumps messages from the WebSocket connection to the manager
func (c *Client) ReadPump() {
	defer func() {
		c.Manager.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(PongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(PongWait))
		c.LastActivity = time.Now()
		return nil
	})

	for {
		select {
		case <-c.Context.Done():
			return
		default:
			_, message, err := c.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Error().Err(err).Str("client_id", c.ID).Msg("WebSocket read error")
				}
				return
			}

			c.LastActivity = time.Now()
			c.handleMessage(message)

			// Update metrics
			c.Manager.metrics.mu.Lock()
			c.Manager.metrics.MessagesReceived++
			c.Manager.metrics.mu.Unlock()
		}
	}
}

// WritePump pumps messages from the manager to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case <-c.Context.Done():
			return

		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				// Channel closed
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current WebSocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

			// Update metrics
			c.Manager.metrics.mu.Lock()
			c.Manager.metrics.MessagesSent++
			c.Manager.metrics.mu.Unlock()

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage handles incoming client messages
func (c *Client) handleMessage(data []byte) {
	var msg ClientMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.sendError("", "Invalid message format")
		log.Error().Err(err).Str("client_id", c.ID).Msg("Failed to parse client message")
		return
	}

	switch msg.Type {
	case MessageTypeSubscribe:
		c.handleSubscribe(msg)
	case MessageTypeUnsubscribe:
		c.handleUnsubscribe(msg)
	case MessageTypePing:
		c.handlePing(msg)
	case MessageTypeTyping:
		c.handleTyping(msg)
	case MessageTypeReadReceipt:
		c.handleReadReceipt(msg)
	default:
		c.sendError(msg.ID, "Unknown message type: "+msg.Type)
	}
}

// handleSubscribe handles subscription requests
func (c *Client) handleSubscribe(msg ClientMessage) {
	if msg.Channel == "" {
		c.sendError(msg.ID, "Channel is required for subscription")
		return
	}

	// Validate channel access (basic validation, enhance as needed)
	if !c.canAccessChannel(msg.Channel) {
		c.sendError(msg.ID, "Access denied to channel: "+msg.Channel)
		log.Warn().
			Str("client_id", c.ID).
			Str("user_id", c.UserID.String()).
			Str("channel", msg.Channel).
			Msg("Unauthorized channel subscription attempt")
		return
	}

	c.SubscriptionsMu.Lock()
	c.Subscriptions[msg.Channel] = true
	c.SubscriptionsMu.Unlock()

	c.sendAck(msg.ID, "Subscribed to "+msg.Channel)

	log.Debug().
		Str("client_id", c.ID).
		Str("channel", msg.Channel).
		Msg("Client subscribed to channel")
}

// handleUnsubscribe handles unsubscription requests
func (c *Client) handleUnsubscribe(msg ClientMessage) {
	if msg.Channel == "" {
		c.sendError(msg.ID, "Channel is required for unsubscription")
		return
	}

	c.SubscriptionsMu.Lock()
	delete(c.Subscriptions, msg.Channel)
	c.SubscriptionsMu.Unlock()

	c.sendAck(msg.ID, "Unsubscribed from "+msg.Channel)

	log.Debug().
		Str("client_id", c.ID).
		Str("channel", msg.Channel).
		Msg("Client unsubscribed from channel")
}

// handlePing handles ping messages
func (c *Client) handlePing(msg ClientMessage) {
	c.sendMessage(ServerMessage{
		Type:      MessageTypePong,
		ID:        msg.ID,
		Timestamp: time.Now(),
	})
}

// handleTyping handles typing indicator messages
func (c *Client) handleTyping(msg ClientMessage) {
	// Typing indicators are ephemeral and don't need persistence
	// Broadcast to other participants in the channel
	if msg.Channel == "" {
		return
	}

	// TODO: Implement typing indicator broadcasting
	// This would typically broadcast to other users in a conversation
	log.Debug().
		Str("client_id", c.ID).
		Str("channel", msg.Channel).
		Msg("Typing indicator received")
}

// handleReadReceipt handles read receipt messages
func (c *Client) handleReadReceipt(msg ClientMessage) {
	// TODO: Implement read receipt handling
	// This would typically update message read status in the database
	log.Debug().
		Str("client_id", c.ID).
		Str("channel", msg.Channel).
		Msg("Read receipt received")
}

// canAccessChannel checks if the client has access to a channel
func (c *Client) canAccessChannel(channel string) bool {
	// Basic validation: users can subscribe to their own user channel
	if channel == Channel.User(c.UserID) {
		return true
	}

	// TODO: Implement more sophisticated authorization:
	// - Check if user is member of the space (for space: channels)
	// - Check if user is participant in conversation (for conv: channels)
	// - Check if user has access to post (for post: channels)
	// For now, allow all subscriptions (to be enhanced)

	return true
}

// sendMessage sends a server message to the client
func (c *Client) sendMessage(msg ServerMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal server message")
		return
	}

	select {
	case c.Send <- data:
	default:
		log.Warn().
			Str("client_id", c.ID).
			Msg("Client send buffer full, dropping message")
	}
}

// sendAck sends an acknowledgment message
func (c *Client) sendAck(msgID, message string) {
	c.sendMessage(ServerMessage{
		Type:      MessageTypeAck,
		ID:        msgID,
		Payload:   map[string]interface{}{"message": message},
		Timestamp: time.Now(),
	})
}

// sendError sends an error message
func (c *Client) sendError(msgID, errorMsg string) {
	c.sendMessage(ServerMessage{
		Type:      MessageTypeError,
		ID:        msgID,
		Error:     errorMsg,
		Timestamp: time.Now(),
	})

	// Update metrics
	c.Manager.metrics.mu.Lock()
	c.Manager.metrics.Errors++
	c.Manager.metrics.LastError = errorMsg
	c.Manager.metrics.LastErrorTime = time.Now()
	c.Manager.metrics.mu.Unlock()
}
