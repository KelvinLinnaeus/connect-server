package websocket

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client connection
type Client struct {
	ID             string                // Unique client connection ID
	UserID         uuid.UUID             // Authenticated user ID
	SpaceID        *uuid.UUID            // Current space context (if any)
	IPAddress      string                // Client IP address
	Conn           *websocket.Conn       // WebSocket connection
	Send           chan []byte           // Buffered channel for outbound messages
	Subscriptions  map[string]bool       // Channels the client is subscribed to
	SubscriptionsMu sync.RWMutex          // Mutex for subscriptions
	Manager        *Manager              // Reference to manager
	Context        context.Context       // Client context
	Cancel         context.CancelFunc    // Cancel function
	LastActivity   time.Time             // Last activity timestamp
	ConnectedAt    time.Time             // Connection timestamp
	Metadata       map[string]string     // Additional client metadata
}

// ClientMessage represents a message from the client
type ClientMessage struct {
	Type    string                 `json:"type"`    // Message type (subscribe, unsubscribe, ping, message, typing, etc.)
	Channel string                 `json:"channel"` // Channel/topic for subscription or event
	Payload map[string]interface{} `json:"payload"` // Message payload
	ID      string                 `json:"id"`      // Optional message ID for tracking
}

// ServerMessage represents a message to the client
type ServerMessage struct {
	Type      string                 `json:"type"`      // Message type (event, ack, error, pong)
	Channel   string                 `json:"channel"`   // Channel/topic
	Payload   map[string]interface{} `json:"payload"`   // Message payload
	ID        string                 `json:"id"`        // Message ID
	Timestamp time.Time              `json:"timestamp"` // Server timestamp
	Error     string                 `json:"error,omitempty"` // Error message if any
}

// Message types
const (
	// Client to server
	MessageTypeSubscribe   = "subscribe"
	MessageTypeUnsubscribe = "unsubscribe"
	MessageTypePing        = "ping"
	MessageTypeMessage     = "message"
	MessageTypeTyping      = "typing"
	MessageTypeReadReceipt = "read"

	// Server to client
	MessageTypeEvent = "event"
	MessageTypeAck   = "ack"
	MessageTypeError = "error"
	MessageTypePong  = "pong"
)

// Connection configuration
const (
	// Time allowed to write a message to the peer
	WriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	PongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	PingPeriod = (PongWait * 9) / 10

	// Heartbeat interval (application-level ping)
	HeartbeatInterval = 30 * time.Second

	// Maximum message size allowed from peer (2MB for production)
	MaxMessageSize = 2 * 1024 * 1024 // 2 MB

	// Send buffer size
	SendBufferSize = 256

	// Time before considering a client idle
	IdleTimeout = 5 * time.Minute

	// Maximum concurrent connections per user
	MaxConnectionsPerUser = 100

	// Maximum concurrent connections per IP address
	MaxConnectionsPerIP = 100
)

// Manager manages all client connections
type Manager struct {
	clients    map[string]*Client      // All connected clients
	clientsMu  sync.RWMutex            // Mutex for clients map
	userIndex  map[uuid.UUID][]*Client // Index of clients by user ID
	indexMu    sync.RWMutex            // Mutex for user index
	ipIndex    map[string][]*Client    // Index of clients by IP address
	ipMu       sync.RWMutex            // Mutex for IP index
	register   chan *Client            // Register client channel
	unregister chan *Client            // Unregister client channel
	broadcast  chan *BroadcastMessage  // Broadcast message channel
	ctx        context.Context         // Manager context
	cancel     context.CancelFunc      // Cancel function
	metrics    *Metrics                // Metrics collector
}

// BroadcastMessage represents a message to broadcast to multiple clients
type BroadcastMessage struct {
	UserIDs []uuid.UUID // Target user IDs (if nil, broadcast to all)
	Channel string      // Channel name
	Message ServerMessage
}

// Metrics tracks WebSocket statistics
type Metrics struct {
	TotalConnections    int64         // Total connections since start
	ActiveConnections   int64         // Currently active connections
	MessagesReceived    int64         // Total messages received
	MessagesSent        int64         // Total messages sent
	Errors              int64         // Total errors
	ConnectionsRejected int64         // Connections rejected due to limits
	LastError           string        // Last error message
	LastErrorTime       time.Time     // Last error time

	// Latency tracking
	TotalLatencyMs      int64         // Sum of all message latencies (ms)
	LatencyCount        int64         // Number of latency measurements

	// Throughput tracking
	StartTime           time.Time     // Metrics collection start time

	mu                  sync.RWMutex  // Mutex for metrics
}

// GetAverageLatencyMs returns average message latency in milliseconds
func (m *Metrics) GetAverageLatencyMs() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.LatencyCount == 0 {
		return 0
	}
	return float64(m.TotalLatencyMs) / float64(m.LatencyCount)
}

// GetMessageThroughput returns messages per second
func (m *Metrics) GetMessageThroughput() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	duration := time.Since(m.StartTime).Seconds()
	if duration == 0 {
		return 0
	}
	return float64(m.MessagesReceived+m.MessagesSent) / duration
}

// ChannelPattern helps construct channel names
type ChannelPattern struct{}

var Channel = &ChannelPattern{}

// User creates a user-specific channel
func (cp *ChannelPattern) User(userID uuid.UUID) string {
	return "user:" + userID.String()
}

// Space creates a space-specific channel
func (cp *ChannelPattern) Space(spaceID uuid.UUID) string {
	return "space:" + spaceID.String()
}

// Conversation creates a conversation-specific channel
func (cp *ChannelPattern) Conversation(convID uuid.UUID) string {
	return "conv:" + convID.String()
}

// Post creates a post-specific channel
func (cp *ChannelPattern) Post(postID uuid.UUID) string {
	return "post:" + postID.String()
}

// Event creates an event-specific channel
func (cp *ChannelPattern) Event(eventID uuid.UUID) string {
	return "event:" + eventID.String()
}
