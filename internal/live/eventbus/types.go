package eventbus

import (
	"context"
	"time"

	"github.com/google/uuid"
)


type Event struct {
	ID        string                 `json:"id"`         
	Type      string                 `json:"type"`       
	Channel   string                 `json:"channel"`    
	Payload   map[string]interface{} `json:"payload"`    
	Timestamp time.Time              `json:"timestamp"`  
	UserID    *uuid.UUID             `json:"user_id"`    
	SpaceID   *uuid.UUID             `json:"space_id"`   
	Metadata  map[string]string      `json:"metadata"`   
}


type EventHandler func(event *Event) error


type EventBus interface {
	
	Publish(ctx context.Context, event *Event) error

	
	
	Subscribe(ctx context.Context, channel string) (<-chan *Event, error)

	
	SubscribePattern(ctx context.Context, pattern string) (<-chan *Event, error)

	
	Unsubscribe(ctx context.Context, channel string) error

	
	Close() error

	
	HealthCheck(ctx context.Context) error
}


const (
	
	EventTypeMessageCreated  = "message.created"
	EventTypeMessageDelivered = "message.delivered"
	EventTypeMessageRead     = "message.read"
	EventTypeMessageDeleted  = "message.deleted"
	EventTypeTypingStarted   = "typing.started"
	EventTypeTypingStopped   = "typing.stopped"

	
	EventTypeNotificationCreated = "notification.created"
	EventTypeNotificationRead    = "notification.read"
	EventTypeNotificationDeleted = "notification.deleted"

	
	EventTypePostCreated  = "post.created"
	EventTypePostUpdated  = "post.updated"
	EventTypePostDeleted  = "post.deleted"
	EventTypePostLiked    = "post.liked"
	EventTypePostUnliked  = "post.unliked"

	
	EventTypeCommentCreated = "comment.created"
	EventTypeCommentUpdated = "comment.updated"
	EventTypeCommentDeleted = "comment.deleted"

	
	EventTypeSpaceUpdated      = "space.updated"
	EventTypeSpaceMemberJoined = "space.member.joined"
	EventTypeSpaceMemberLeft   = "space.member.left"

	
	EventTypeLessonPublished  = "lesson.published"
	EventTypeSessionStarted   = "session.started"
	EventTypeAssignmentGraded = "assignment.graded"

	
	EventTypeEventUpdated    = "event.updated"
	EventTypeEventCancelled  = "event.cancelled"
	EventTypeEventRSVPAdded  = "event.rsvp.added"

	
	EventTypeUserOnline  = "user.online"
	EventTypeUserOffline = "user.offline"
	EventTypeUserIdle    = "user.idle"
)


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


func (e *Event) WithUserID(userID uuid.UUID) *Event {
	e.UserID = &userID
	return e
}


func (e *Event) WithSpaceID(spaceID uuid.UUID) *Event {
	e.SpaceID = &spaceID
	return e
}


func (e *Event) WithMetadata(key, value string) *Event {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
	return e
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
