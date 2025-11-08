package eventbus

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Event represents a real-time event that can be published and subscribed to
type Event struct {
	ID        string                 `json:"id"`         // Unique event ID
	Type      string                 `json:"type"`       // Event type (e.g., "message.created", "notification.created")
	Channel   string                 `json:"channel"`    // Channel/topic (e.g., "conv:uuid", "space:uuid", "user:uuid")
	Payload   map[string]interface{} `json:"payload"`    // Event data
	Timestamp time.Time              `json:"timestamp"`  // When the event occurred
	UserID    *uuid.UUID             `json:"user_id"`    // User who triggered the event (if applicable)
	SpaceID   *uuid.UUID             `json:"space_id"`   // Space context (if applicable)
	Metadata  map[string]string      `json:"metadata"`   // Additional metadata
}

// EventHandler is a function that handles received events
type EventHandler func(event *Event) error

// EventBus defines the interface for publishing and subscribing to events
type EventBus interface {
	// Publish publishes an event to a channel
	Publish(ctx context.Context, event *Event) error

	// Subscribe subscribes to a channel and receives events
	// Returns a channel that receives events and an error if subscription fails
	Subscribe(ctx context.Context, channel string) (<-chan *Event, error)

	// SubscribePattern subscribes to channels matching a pattern (e.g., "user:*")
	SubscribePattern(ctx context.Context, pattern string) (<-chan *Event, error)

	// Unsubscribe unsubscribes from a channel
	Unsubscribe(ctx context.Context, channel string) error

	// Close closes the event bus and cleans up resources
	Close() error

	// HealthCheck checks if the event bus is healthy
	HealthCheck(ctx context.Context) error
}

// EventType constants for common event types
const (
	// Messaging events
	EventTypeMessageCreated  = "message.created"
	EventTypeMessageDelivered = "message.delivered"
	EventTypeMessageRead     = "message.read"
	EventTypeMessageDeleted  = "message.deleted"
	EventTypeTypingStarted   = "typing.started"
	EventTypeTypingStopped   = "typing.stopped"

	// Notification events
	EventTypeNotificationCreated = "notification.created"
	EventTypeNotificationRead    = "notification.read"
	EventTypeNotificationDeleted = "notification.deleted"

	// Post events
	EventTypePostCreated  = "post.created"
	EventTypePostUpdated  = "post.updated"
	EventTypePostDeleted  = "post.deleted"
	EventTypePostLiked    = "post.liked"
	EventTypePostUnliked  = "post.unliked"

	// Comment events
	EventTypeCommentCreated = "comment.created"
	EventTypeCommentUpdated = "comment.updated"
	EventTypeCommentDeleted = "comment.deleted"

	// Space events
	EventTypeSpaceUpdated      = "space.updated"
	EventTypeSpaceMemberJoined = "space.member.joined"
	EventTypeSpaceMemberLeft   = "space.member.left"

	// Learning events
	EventTypeLessonPublished  = "lesson.published"
	EventTypeSessionStarted   = "session.started"
	EventTypeAssignmentGraded = "assignment.graded"

	// Event (calendar) events
	EventTypeEventUpdated    = "event.updated"
	EventTypeEventCancelled  = "event.cancelled"
	EventTypeEventRSVPAdded  = "event.rsvp.added"

	// Presence events
	EventTypeUserOnline  = "user.online"
	EventTypeUserOffline = "user.offline"
	EventTypeUserIdle    = "user.idle"
)

// NewEvent creates a new event with defaults
func NewEvent(eventType, channel string, payload map[string]interface{}) *Event {
	return &Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Channel:   channel,
		Payload:   payload,
		Timestamp: time.Now(),
		Metadata:  make(map[string]string),
	}
}

// WithUserID sets the user ID for the event
func (e *Event) WithUserID(userID uuid.UUID) *Event {
	e.UserID = &userID
	return e
}

// WithSpaceID sets the space ID for the event
func (e *Event) WithSpaceID(spaceID uuid.UUID) *Event {
	e.SpaceID = &spaceID
	return e
}

// WithMetadata adds metadata to the event
func (e *Event) WithMetadata(key, value string) *Event {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
	return e
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
