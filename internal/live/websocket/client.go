package websocket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)


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

			
			c.Manager.metrics.mu.Lock()
			c.Manager.metrics.MessagesReceived++
			c.Manager.metrics.mu.Unlock()
		}
	}
}


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
				
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

			
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


func (c *Client) handleSubscribe(msg ClientMessage) {
	if msg.Channel == "" {
		c.sendError(msg.ID, "Channel is required for subscription")
		return
	}

	
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


func (c *Client) handlePing(msg ClientMessage) {
	c.sendMessage(ServerMessage{
		Type:      MessageTypePong,
		ID:        msg.ID,
		Timestamp: time.Now(),
	})
}


func (c *Client) handleTyping(msg ClientMessage) {
	
	
	if msg.Channel == "" {
		return
	}

	
	
	log.Debug().
		Str("client_id", c.ID).
		Str("channel", msg.Channel).
		Msg("Typing indicator received")
}


func (c *Client) handleReadReceipt(msg ClientMessage) {
	
	
	log.Debug().
		Str("client_id", c.ID).
		Str("channel", msg.Channel).
		Msg("Read receipt received")
}


func (c *Client) canAccessChannel(channel string) bool {
	
	if channel == Channel.User(c.UserID) {
		return true
	}

	
	
	
	
	

	return true
}


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


func (c *Client) sendAck(msgID, message string) {
	c.sendMessage(ServerMessage{
		Type:      MessageTypeAck,
		ID:        msgID,
		Payload:   map[string]interface{}{"message": message},
		Timestamp: time.Now(),
	})
}


func (c *Client) sendError(msgID, errorMsg string) {
	c.sendMessage(ServerMessage{
		Type:      MessageTypeError,
		ID:        msgID,
		Error:     errorMsg,
		Timestamp: time.Now(),
	})

	
	c.Manager.metrics.mu.Lock()
	c.Manager.metrics.Errors++
	c.Manager.metrics.LastError = errorMsg
	c.Manager.metrics.LastErrorTime = time.Now()
	c.Manager.metrics.mu.Unlock()
}
