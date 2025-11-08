package websocket

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/connect-univyn/connect-server/internal/live/eventbus"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)


func NewManager(ctx context.Context, bus eventbus.EventBus) *Manager {
	managerCtx, cancel := context.WithCancel(ctx)

	m := &Manager{
		clients:    make(map[string]*Client),
		userIndex:  make(map[uuid.UUID][]*Client),
		ipIndex:    make(map[string][]*Client),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		broadcast:  make(chan *BroadcastMessage, 1024),
		ctx:        managerCtx,
		cancel:     cancel,
		metrics:    &Metrics{StartTime: time.Now()},
	}

	
	go m.run()
	go m.listenToEventBus(bus)

	log.Info().
		Dur("heartbeat", HeartbeatInterval).
		Int("max_msg_size_mb", MaxMessageSize/(1024*1024)).
		Int("max_conn_per_user", MaxConnectionsPerUser).
		Int("max_conn_per_ip", MaxConnectionsPerIP).
		Msg("WebSocket hub initialized")

	return m
}


func (m *Manager) run() {
	ticker := time.NewTicker(PingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			log.Info().Msg("WebSocket manager shutting down")
			return

		case client := <-m.register:
			m.registerClient(client)

		case client := <-m.unregister:
			m.unregisterClient(client)

		case broadcast := <-m.broadcast:
			m.broadcastMessage(broadcast)

		case <-ticker.C:
			
			m.cleanupIdleClients()
		}
	}
}


func (m *Manager) registerClient(client *Client) {
	
	m.indexMu.RLock()
	userConnCount := len(m.userIndex[client.UserID])
	m.indexMu.RUnlock()

	if userConnCount >= MaxConnectionsPerUser {
		m.metrics.mu.Lock()
		m.metrics.ConnectionsRejected++
		m.metrics.mu.Unlock()

		log.Warn().
			Str("user_id", client.UserID.String()).
			Int("current_connections", userConnCount).
			Int("max_allowed", MaxConnectionsPerUser).
			Msg("Connection rejected: max connections per user exceeded")

		
		client.Conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Maximum connections per user exceeded"),
			time.Now().Add(WriteWait),
		)
		client.Conn.Close()
		return
	}

	m.ipMu.RLock()
	ipConnCount := len(m.ipIndex[client.IPAddress])
	m.ipMu.RUnlock()

	if ipConnCount >= MaxConnectionsPerIP {
		m.metrics.mu.Lock()
		m.metrics.ConnectionsRejected++
		m.metrics.mu.Unlock()

		log.Warn().
			Str("ip_address", client.IPAddress).
			Int("current_connections", ipConnCount).
			Int("max_allowed", MaxConnectionsPerIP).
			Msg("Connection rejected: max connections per IP exceeded")

		
		client.Conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Maximum connections per IP exceeded"),
			time.Now().Add(WriteWait),
		)
		client.Conn.Close()
		return
	}

	
	m.clientsMu.Lock()
	m.clients[client.ID] = client
	m.clientsMu.Unlock()

	
	m.indexMu.Lock()
	m.userIndex[client.UserID] = append(m.userIndex[client.UserID], client)
	m.indexMu.Unlock()

	
	m.ipMu.Lock()
	m.ipIndex[client.IPAddress] = append(m.ipIndex[client.IPAddress], client)
	m.ipMu.Unlock()

	
	m.metrics.mu.Lock()
	m.metrics.TotalConnections++
	m.metrics.ActiveConnections++
	m.metrics.mu.Unlock()

	log.Info().
		Str("client_id", client.ID).
		Str("user_id", client.UserID.String()).
		Str("ip_address", client.IPAddress).
		Int64("active_connections", m.metrics.ActiveConnections).
		Msg("Client registered")
}


func (m *Manager) unregisterClient(client *Client) {
	m.clientsMu.Lock()
	if _, exists := m.clients[client.ID]; exists {
		delete(m.clients, client.ID)
		close(client.Send)
	}
	m.clientsMu.Unlock()

	
	m.indexMu.Lock()
	if clients, exists := m.userIndex[client.UserID]; exists {
		for i, c := range clients {
			if c.ID == client.ID {
				m.userIndex[client.UserID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		
		if len(m.userIndex[client.UserID]) == 0 {
			delete(m.userIndex, client.UserID)
		}
	}
	m.indexMu.Unlock()

	
	m.ipMu.Lock()
	if clients, exists := m.ipIndex[client.IPAddress]; exists {
		for i, c := range clients {
			if c.ID == client.ID {
				m.ipIndex[client.IPAddress] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		
		if len(m.ipIndex[client.IPAddress]) == 0 {
			delete(m.ipIndex, client.IPAddress)
		}
	}
	m.ipMu.Unlock()

	
	m.metrics.mu.Lock()
	m.metrics.ActiveConnections--
	m.metrics.mu.Unlock()

	
	duration := time.Since(client.ConnectedAt)

	log.Info().
		Str("client_id", client.ID).
		Str("user_id", client.UserID.String()).
		Str("ip_address", client.IPAddress).
		Dur("session_duration", duration).
		Int64("active_connections", m.metrics.ActiveConnections).
		Msg("Client unregistered")
}


func (m *Manager) broadcastMessage(broadcast *BroadcastMessage) {
	var targetClients []*Client

	if broadcast.UserIDs == nil {
		
		m.clientsMu.RLock()
		targetClients = make([]*Client, 0, len(m.clients))
		for _, client := range m.clients {
			targetClients = append(targetClients, client)
		}
		m.clientsMu.RUnlock()
	} else {
		
		m.indexMu.RLock()
		for _, userID := range broadcast.UserIDs {
			if clients, exists := m.userIndex[userID]; exists {
				targetClients = append(targetClients, clients...)
			}
		}
		m.indexMu.RUnlock()
	}

	
	data, err := json.Marshal(broadcast.Message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal broadcast message")
		return
	}

	sent := 0
	for _, client := range targetClients {
		
		client.SubscriptionsMu.RLock()
		subscribed := client.Subscriptions[broadcast.Channel] || broadcast.Channel == ""
		client.SubscriptionsMu.RUnlock()

		if subscribed {
			select {
			case client.Send <- data:
				sent++
			default:
				log.Warn().
					Str("client_id", client.ID).
					Str("user_id", client.UserID.String()).
					Msg("Client send buffer full, skipping message")
			}
		}
	}

	log.Debug().
		Int("target_count", len(targetClients)).
		Int("sent_count", sent).
		Str("channel", broadcast.Channel).
		Msg("Broadcast completed")
}


func (m *Manager) listenToEventBus(bus eventbus.EventBus) {
	
	
	eventChan, err := bus.SubscribePattern(m.ctx, "*")
	if err != nil {
		log.Error().Err(err).Msg("Failed to subscribe to event bus")
		return
	}

	log.Info().Msg("WebSocket manager listening to event bus")

	for {
		select {
		case <-m.ctx.Done():
			return
		case event, ok := <-eventChan:
			if !ok {
				log.Warn().Msg("Event bus channel closed")
				return
			}

			
			serverMsg := ServerMessage{
				Type:      MessageTypeEvent,
				Channel:   event.Channel,
				Payload:   event.Payload,
				ID:        event.ID,
				Timestamp: event.Timestamp,
			}

			
			var userIDs []uuid.UUID
			if event.UserID != nil {
				userIDs = []uuid.UUID{*event.UserID}
			}

			m.broadcast <- &BroadcastMessage{
				UserIDs: userIDs,
				Channel: event.Channel,
				Message: serverMsg,
			}
		}
	}
}


func (m *Manager) cleanupIdleClients() {
	now := time.Now()
	var idleClients []*Client

	m.clientsMu.RLock()
	for _, client := range m.clients {
		if now.Sub(client.LastActivity) > IdleTimeout {
			idleClients = append(idleClients, client)
		}
	}
	m.clientsMu.RUnlock()

	for _, client := range idleClients {
		log.Info().
			Str("client_id", client.ID).
			Str("user_id", client.UserID.String()).
			Dur("idle_duration", now.Sub(client.LastActivity)).
			Msg("Removing idle client")
		m.unregister <- client
		client.Cancel()
	}
}


func (m *Manager) Register(client *Client) {
	m.register <- client
}


func (m *Manager) Unregister(client *Client) {
	m.unregister <- client
}


func (m *Manager) Broadcast(userIDs []uuid.UUID, channel string, message ServerMessage) {
	m.broadcast <- &BroadcastMessage{
		UserIDs: userIDs,
		Channel: channel,
		Message: message,
	}
}


func (m *Manager) GetActiveConnections() int64 {
	m.metrics.mu.RLock()
	defer m.metrics.mu.RUnlock()
	return m.metrics.ActiveConnections
}


func (m *Manager) GetMetrics() Metrics {
	m.metrics.mu.RLock()
	defer m.metrics.mu.RUnlock()
	return *m.metrics
}


func (m *Manager) GetUserConnections(userID uuid.UUID) []*Client {
	m.indexMu.RLock()
	defer m.indexMu.RUnlock()

	if clients, exists := m.userIndex[userID]; exists {
		
		result := make([]*Client, len(clients))
		copy(result, clients)
		return result
	}

	return nil
}


func (m *Manager) IsUserOnline(userID uuid.UUID) bool {
	m.indexMu.RLock()
	defer m.indexMu.RUnlock()

	clients, exists := m.userIndex[userID]
	return exists && len(clients) > 0
}


func (m *Manager) Shutdown() error {
	log.Info().Msg("Shutting down WebSocket manager")

	
	m.cancel()

	
	m.clientsMu.RLock()
	clients := make([]*Client, 0, len(m.clients))
	for _, client := range m.clients {
		clients = append(clients, client)
	}
	m.clientsMu.RUnlock()

	var wg sync.WaitGroup
	for _, client := range clients {
		wg.Add(1)
		go func(c *Client) {
			defer wg.Done()
			c.Cancel()
			c.Conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseGoingAway, "Server shutting down"),
				time.Now().Add(WriteWait),
			)
			c.Conn.Close()
		}(client)
	}

	
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info().Msg("All clients disconnected")
	case <-time.After(10 * time.Second):
		log.Warn().Msg("Timeout waiting for clients to disconnect")
	}

	return nil
}
