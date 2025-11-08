package live

import (
	"context"

	"github.com/connect-univyn/connect-server/internal/live/eventbus"
	"github.com/connect-univyn/connect-server/internal/live/websocket"
	"github.com/google/uuid"
)


type Service struct {
	bus       eventbus.EventBus
	wsManager *websocket.Manager
}


func NewService(bus eventbus.EventBus) *Service {
	return &Service{bus: bus}
}


func (s *Service) SetWebSocketManager(manager *websocket.Manager) {
	s.wsManager = manager
}


func (s *Service) GetWebSocketMetrics() websocket.Metrics {
	if s.wsManager == nil {
		return websocket.Metrics{}
	}
	return s.wsManager.GetMetrics()
}


func (s *Service) GetBrokerMetrics() *eventbus.BrokerMetrics {
	if redisBroker, ok := s.bus.(*eventbus.RedisBroker); ok {
		metrics := redisBroker.GetMetrics()
		return &metrics
	}
	return nil
}


func (s *Service) PublishMessageCreated(ctx context.Context, conversationID, senderID uuid.UUID, message map[string]interface{}) error {
	event := eventbus.NewEvent(
		eventbus.EventTypeMessageCreated,
		eventbus.Channel.Conversation(conversationID),
		message,
	).WithUserID(senderID)

	return s.bus.Publish(ctx, event)
}


func (s *Service) PublishMessageDelivered(ctx context.Context, conversationID, messageID, recipientID uuid.UUID) error {
	event := eventbus.NewEvent(
		eventbus.EventTypeMessageDelivered,
		eventbus.Channel.Conversation(conversationID),
		map[string]interface{}{
			"message_id":   messageID.String(),
			"recipient_id": recipientID.String(),
		},
	).WithUserID(recipientID)

	return s.bus.Publish(ctx, event)
}


func (s *Service) PublishMessageRead(ctx context.Context, conversationID uuid.UUID, messageIDs []uuid.UUID, userID uuid.UUID) error {
	msgIDStrs := make([]string, len(messageIDs))
	for i, id := range messageIDs {
		msgIDStrs[i] = id.String()
	}

	event := eventbus.NewEvent(
		eventbus.EventTypeMessageRead,
		eventbus.Channel.Conversation(conversationID),
		map[string]interface{}{
			"message_ids": msgIDStrs,
			"user_id":     userID.String(),
		},
	).WithUserID(userID)

	return s.bus.Publish(ctx, event)
}


func (s *Service) PublishTypingStarted(ctx context.Context, conversationID, userID uuid.UUID, username string) error {
	event := eventbus.NewEvent(
		eventbus.EventTypeTypingStarted,
		eventbus.Channel.Conversation(conversationID),
		map[string]interface{}{
			"user_id":  userID.String(),
			"username": username,
		},
	).WithUserID(userID)

	return s.bus.Publish(ctx, event)
}


func (s *Service) PublishTypingStopped(ctx context.Context, conversationID, userID uuid.UUID, username string) error {
	event := eventbus.NewEvent(
		eventbus.EventTypeTypingStopped,
		eventbus.Channel.Conversation(conversationID),
		map[string]interface{}{
			"user_id":  userID.String(),
			"username": username,
		},
	).WithUserID(userID)

	return s.bus.Publish(ctx, event)
}


func (s *Service) PublishNotificationCreated(ctx context.Context, userID uuid.UUID, notification map[string]interface{}) error {
	event := eventbus.NewEvent(
		eventbus.EventTypeNotificationCreated,
		eventbus.Channel.User(userID),
		notification,
	).WithUserID(userID)

	return s.bus.Publish(ctx, event)
}


func (s *Service) PublishPostCreated(ctx context.Context, spaceID, authorID uuid.UUID, post map[string]interface{}) error {
	event := eventbus.NewEvent(
		eventbus.EventTypePostCreated,
		eventbus.Channel.Space(spaceID),
		post,
	).WithUserID(authorID).WithSpaceID(spaceID)

	return s.bus.Publish(ctx, event)
}


func (s *Service) PublishPostUpdated(ctx context.Context, spaceID, postID, authorID uuid.UUID, updates map[string]interface{}) error {
	event := eventbus.NewEvent(
		eventbus.EventTypePostUpdated,
		eventbus.Channel.Post(postID),
		updates,
	).WithUserID(authorID).WithSpaceID(spaceID)

	
	spaceEvent := eventbus.NewEvent(
		eventbus.EventTypePostUpdated,
		eventbus.Channel.Space(spaceID),
		updates,
	).WithUserID(authorID).WithSpaceID(spaceID)

	if err := s.bus.Publish(ctx, event); err != nil {
		return err
	}
	return s.bus.Publish(ctx, spaceEvent)
}


func (s *Service) PublishPostLiked(ctx context.Context, postID, userID, spaceID uuid.UUID, likeCount int) error {
	event := eventbus.NewEvent(
		eventbus.EventTypePostLiked,
		eventbus.Channel.Post(postID),
		map[string]interface{}{
			"post_id":    postID.String(),
			"user_id":    userID.String(),
			"like_count": likeCount,
		},
	).WithUserID(userID).WithSpaceID(spaceID)

	return s.bus.Publish(ctx, event)
}


func (s *Service) PublishCommentCreated(ctx context.Context, postID, authorID uuid.UUID, comment map[string]interface{}) error {
	event := eventbus.NewEvent(
		eventbus.EventTypeCommentCreated,
		eventbus.Channel.Post(postID),
		comment,
	).WithUserID(authorID)

	return s.bus.Publish(ctx, event)
}


func (s *Service) PublishSpaceMemberJoined(ctx context.Context, spaceID, userID uuid.UUID, memberInfo map[string]interface{}) error {
	event := eventbus.NewEvent(
		eventbus.EventTypeSpaceMemberJoined,
		eventbus.Channel.Space(spaceID),
		memberInfo,
	).WithUserID(userID).WithSpaceID(spaceID)

	return s.bus.Publish(ctx, event)
}


func (s *Service) PublishSpaceMemberLeft(ctx context.Context, spaceID, userID uuid.UUID) error {
	event := eventbus.NewEvent(
		eventbus.EventTypeSpaceMemberLeft,
		eventbus.Channel.Space(spaceID),
		map[string]interface{}{
			"user_id": userID.String(),
		},
	).WithUserID(userID).WithSpaceID(spaceID)

	return s.bus.Publish(ctx, event)
}


func (s *Service) PublishUserOnline(ctx context.Context, userID uuid.UUID, metadata map[string]string) error {
	payload := map[string]interface{}{
		"user_id": userID.String(),
		"status":  "online",
	}
	if metadata != nil {
		payload["metadata"] = metadata
	}

	event := eventbus.NewEvent(
		eventbus.EventTypeUserOnline,
		eventbus.Channel.User(userID),
		payload,
	).WithUserID(userID)

	return s.bus.Publish(ctx, event)
}


func (s *Service) PublishUserOffline(ctx context.Context, userID uuid.UUID) error {
	event := eventbus.NewEvent(
		eventbus.EventTypeUserOffline,
		eventbus.Channel.User(userID),
		map[string]interface{}{
			"user_id": userID.String(),
			"status":  "offline",
		},
	).WithUserID(userID)

	return s.bus.Publish(ctx, event)
}


func (s *Service) PublishLessonPublished(ctx context.Context, spaceID uuid.UUID, lesson map[string]interface{}) error {
	event := eventbus.NewEvent(
		eventbus.EventTypeLessonPublished,
		eventbus.Channel.Space(spaceID),
		lesson,
	).WithSpaceID(spaceID)

	return s.bus.Publish(ctx, event)
}


func (s *Service) PublishEventUpdated(ctx context.Context, eventID, spaceID uuid.UUID, updates map[string]interface{}) error {
	event := eventbus.NewEvent(
		eventbus.EventTypeEventUpdated,
		eventbus.Channel.Event(eventID),
		updates,
	).WithSpaceID(spaceID)

	
	spaceEvent := eventbus.NewEvent(
		eventbus.EventTypeEventUpdated,
		eventbus.Channel.Space(spaceID),
		updates,
	).WithSpaceID(spaceID)

	if err := s.bus.Publish(ctx, event); err != nil {
		return err
	}
	return s.bus.Publish(ctx, spaceEvent)
}


func (s *Service) PublishEvent(ctx context.Context, eventType, channel string, payload map[string]interface{}, userID *uuid.UUID, spaceID *uuid.UUID) error {
	event := eventbus.NewEvent(eventType, channel, payload)
	if userID != nil {
		event = event.WithUserID(*userID)
	}
	if spaceID != nil {
		event = event.WithSpaceID(*spaceID)
	}

	return s.bus.Publish(ctx, event)
}
