package websocket

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)


type Client struct {
	ID             string                
	UserID         uuid.UUID             
	SpaceID        *uuid.UUID            
	IPAddress      string                
	Conn           *websocket.Conn       
	Send           chan []byte           
	Subscriptions  map[string]bool       
	SubscriptionsMu sync.RWMutex          
	Manager        *Manager              
	Context        context.Context       
	Cancel         context.CancelFunc    
	LastActivity   time.Time             
	ConnectedAt    time.Time             
	Metadata       map[string]string     
}


type ClientMessage struct {
	Type    string                 `json:"type"`    
	Channel string                 `json:"channel"` 
	Payload map[string]interface{} `json:"payload"` 
	ID      string                 `json:"id"`      
}


type ServerMessage struct {
	Type      string                 `json:"type"`      
	Channel   string                 `json:"channel"`   
	Payload   map[string]interface{} `json:"payload"`   
	ID        string                 `json:"id"`        
	Timestamp time.Time              `json:"timestamp"` 
	Error     string                 `json:"error,omitempty"` 
}


const (
	
	MessageTypeSubscribe   = "subscribe"
	MessageTypeUnsubscribe = "unsubscribe"
	MessageTypePing        = "ping"
	MessageTypeMessage     = "message"
	MessageTypeTyping      = "typing"
	MessageTypeReadReceipt = "read"

	
	MessageTypeEvent = "event"
	MessageTypeAck   = "ack"
	MessageTypeError = "error"
	MessageTypePong  = "pong"
)


const (
	
	WriteWait = 10 * time.Second

	
	PongWait = 60 * time.Second

	
	PingPeriod = (PongWait * 9) / 10

	
	HeartbeatInterval = 30 * time.Second

	
	MaxMessageSize = 2 * 1024 * 1024 

	
	SendBufferSize = 256

	
	IdleTimeout = 5 * time.Minute

	
	MaxConnectionsPerUser = 100

	
	MaxConnectionsPerIP = 100
)


type Manager struct {
	clients    map[string]*Client      
	clientsMu  sync.RWMutex            
	userIndex  map[uuid.UUID][]*Client 
	indexMu    sync.RWMutex            
	ipIndex    map[string][]*Client    
	ipMu       sync.RWMutex            
	register   chan *Client            
	unregister chan *Client            
	broadcast  chan *BroadcastMessage  
	ctx        context.Context         
	cancel     context.CancelFunc      
	metrics    *Metrics                
}


type BroadcastMessage struct {
	UserIDs []uuid.UUID 
	Channel string      
	Message ServerMessage
}


type Metrics struct {
	TotalConnections    int64         
	ActiveConnections   int64         
	MessagesReceived    int64         
	MessagesSent        int64         
	Errors              int64         
	ConnectionsRejected int64         
	LastError           string        
	LastErrorTime       time.Time     

	
	TotalLatencyMs      int64         
	LatencyCount        int64         

	
	StartTime           time.Time     

	mu                  sync.RWMutex  
}


func (m *Metrics) GetAverageLatencyMs() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.LatencyCount == 0 {
		return 0
	}
	return float64(m.TotalLatencyMs) / float64(m.LatencyCount)
}


func (m *Metrics) GetMessageThroughput() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	duration := time.Since(m.StartTime).Seconds()
	if duration == 0 {
		return 0
	}
	return float64(m.MessagesReceived+m.MessagesSent) / duration
}


type ChannelPattern struct{}

var Channel = &ChannelPattern{}


func (cp *ChannelPattern) User(userID uuid.UUID) string {
	return "user:" + userID.String()
}


func (cp *ChannelPattern) Space(spaceID uuid.UUID) string {
	return "space:" + spaceID.String()
}


func (cp *ChannelPattern) Conversation(convID uuid.UUID) string {
	return "conv:" + convID.String()
}


func (cp *ChannelPattern) Post(postID uuid.UUID) string {
	return "post:" + postID.String()
}


func (cp *ChannelPattern) Event(eventID uuid.UUID) string {
	return "event:" + eventID.String()
}
