package messaging

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	"github.com/connect-univyn/connect_server/internal/live"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/sqlc-dev/pqtype"
)

// Service handles messaging business logic
type Service struct {
	store       db.Store
	liveService *live.Service
}

// NewService creates a new messaging service
func NewService(store db.Store, liveService *live.Service) *Service {
	return &Service{
		store:       store,
		liveService: liveService,
	}
}

// CreateConversation creates a new conversation
func (s *Service) CreateConversation(ctx context.Context, req CreateConversationRequest) (*ConversationResponse, error) {
	var name, avatar, description, conversationType sql.NullString
	var settings pqtype.NullRawMessage
	
	if req.Name != nil {
		name = sql.NullString{String: *req.Name, Valid: true}
	}
	if req.Avatar != nil {
		avatar = sql.NullString{String: *req.Avatar, Valid: true}
	}
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}
	if req.ConversationType != "" {
		conversationType = sql.NullString{String: req.ConversationType, Valid: true}
	}
	if req.Settings != nil {
		settings = *req.Settings
	}
	
	conversation, err := s.store.CreateConversation(ctx, db.CreateConversationParams{
		SpaceID:          req.SpaceID,
		Name:             name,
		Avatar:           avatar,
		Description:      description,
		ConversationType: conversationType,
		Settings:         settings,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}
	
	// Add participants
	if len(req.ParticipantIDs) > 0 {
		err = s.store.AddConversationParticipants(ctx, db.AddConversationParticipantsParams{
			ConversationID: conversation.ID,
			Column2:        req.ParticipantIDs,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add participants: %w", err)
		}
	}
	
	return s.toConversationResponse(conversation, 0), nil
}

// GetConversationByID gets a conversation by ID
func (s *Service) GetConversationByID(ctx context.Context, conversationID, userID uuid.UUID) (*ConversationResponse, error) {
	conversation, err := s.store.GetConversationByID(ctx, db.GetConversationByIDParams{
		RecipientID: uuid.NullUUID{UUID: userID, Valid: true},
		ID:          conversationID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	
	return s.toConversationDetailResponse(conversation), nil
}

// GetUserConversations gets all conversations for a user
func (s *Service) GetUserConversations(ctx context.Context, userID uuid.UUID) ([]ConversationDetailResponse, error) {
	conversations, err := s.store.GetUserConversations(ctx, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get user conversations: %w", err)
	}
	
	return s.toUserConversationResponses(conversations), nil
}

// GetOrCreateDirectConversation gets existing or creates new direct conversation
func (s *Service) GetOrCreateDirectConversation(ctx context.Context, spaceID, user1ID, user2ID uuid.UUID) (uuid.UUID, error) {
	conversationID, err := s.store.GetOrCreateDirectConversation(ctx, db.GetOrCreateDirectConversationParams{
		SpaceID:  spaceID,
		UserID:   user1ID,
		UserID_2: user2ID,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get or create direct conversation: %w", err)
	}
	
	// If conversation was just created, add both participants
	// Note: The SQLC query creates the conversation but doesn't add participants
	// We need to add them here
	
	return conversationID, nil
}

// SendMessage sends a message in a conversation
func (s *Service) SendMessage(ctx context.Context, req SendMessageRequest) (*MessageResponse, error) {
	var recipientID, replyToID uuid.NullUUID
	var content, messageType sql.NullString
	var attachments pqtype.NullRawMessage
	
	if req.RecipientID != nil {
		recipientID = uuid.NullUUID{UUID: *req.RecipientID, Valid: true}
	}
	if req.ReplyToID != nil {
		replyToID = uuid.NullUUID{UUID: *req.ReplyToID, Valid: true}
	}
	if req.Content != "" {
		content = sql.NullString{String: req.Content, Valid: true}
	}
	if req.MessageType != "" {
		messageType = sql.NullString{String: req.MessageType, Valid: true}
	} else {
		messageType = sql.NullString{String: "text", Valid: true}
	}
	if req.Attachments != nil {
		attachments = *req.Attachments
	}
	
	message, err := s.store.SendMessage(ctx, db.SendMessageParams{
		ConversationID: req.ConversationID,
		SenderID:       req.SenderID,
		RecipientID:    recipientID,
		Content:        content,
		Attachments:    attachments,
		MessageType:    messageType,
		ReplyToID:      replyToID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}
	
	// Update conversation last message asynchronously
	go s.store.UpdateConversationLastMessage(context.Background(), db.UpdateConversationLastMessageParams{
		LastMessageID: uuid.NullUUID{UUID: message.ID, Valid: true},
		ID:            req.ConversationID,
	})
	
	// Get sender info for response
	messageDetail, err := s.store.GetMessageByID(ctx, message.ID)
	if err != nil {
		// Return basic message if we can't get details
		return s.toBasicMessageResponse(message), nil
	}

	response := s.toMessageResponse(messageDetail)

	// Publish real-time event for message creation
	if s.liveService != nil {
		messagePayload := map[string]interface{}{
			"id":              message.ID.String(),
			"conversation_id": message.ConversationID.String(),
			"sender_id":       message.SenderID.String(),
			"content":         response.Content,
			"message_type":    response.MessageType,
			"created_at":      message.CreatedAt.Time.Unix(),
			"sender_username": response.SenderUsername,
			"sender_fullname": response.SenderFullName,
		}
		if response.SenderAvatar != nil {
			messagePayload["sender_avatar"] = *response.SenderAvatar
		}
		if response.Attachments != nil {
			messagePayload["attachments"] = response.Attachments
		}

		if err := s.liveService.PublishMessageCreated(ctx, req.ConversationID, req.SenderID, messagePayload); err != nil {
			log.Error().Err(err).Msg("Failed to publish message.created event")
		}
	}

	return response, nil
}

// GetConversationMessages gets messages for a conversation
func (s *Service) GetConversationMessages(ctx context.Context, params GetConversationMessagesParams) ([]MessageDetailResponse, error) {
	offset := (params.Page - 1) * params.Limit
	
	messages, err := s.store.GetConversationMessages(ctx, db.GetConversationMessagesParams{
		ConversationID: params.ConversationID,
		Limit:          params.Limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation messages: %w", err)
	}
	
	return s.toMessageDetailResponses(messages), nil
}

// GetMessageByID gets a message by ID
func (s *Service) GetMessageByID(ctx context.Context, messageID uuid.UUID) (*MessageResponse, error) {
	message, err := s.store.GetMessageByID(ctx, messageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	
	return s.toMessageResponse(message), nil
}

// DeleteMessage deletes a message (soft delete)
func (s *Service) DeleteMessage(ctx context.Context, messageID, senderID uuid.UUID) error {
	err := s.store.DeleteMessage(ctx, db.DeleteMessageParams{
		ID:       messageID,
		SenderID: senderID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	
	return nil
}

// MarkMessagesAsRead marks all unread messages in a conversation as read
func (s *Service) MarkMessagesAsRead(ctx context.Context, conversationID, userID uuid.UUID) error {
	// Get unread message IDs before marking as read
	messages, err := s.store.GetConversationMessages(ctx, db.GetConversationMessagesParams{
		ConversationID: conversationID,
		Limit:          100, // Get recent unread messages
		Offset:         0,
	})
	if err != nil {
		return fmt.Errorf("failed to get messages: %w", err)
	}

	// Filter to only unread messages for this user
	var messageIDs []uuid.UUID
	for _, msg := range messages {
		if msg.RecipientID.Valid && msg.RecipientID.UUID == userID && msg.IsRead.Valid && !msg.IsRead.Bool {
			messageIDs = append(messageIDs, msg.ID)
		}
	}

	err = s.store.MarkMessagesAsRead(ctx, db.MarkMessagesAsReadParams{
		ConversationID: conversationID,
		RecipientID:    uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}

	// Publish real-time event for message read
	if s.liveService != nil && len(messageIDs) > 0 {
		if err := s.liveService.PublishMessageRead(ctx, conversationID, messageIDs, userID); err != nil {
			log.Error().Err(err).Msg("Failed to publish message.read event")
		}
	}

	return nil
}

// GetUnreadMessageCount gets unread message count for a conversation
func (s *Service) GetUnreadMessageCount(ctx context.Context, conversationID, userID uuid.UUID) (int64, error) {
	count, err := s.store.GetUnreadMessageCount(ctx, db.GetUnreadMessageCountParams{
		ConversationID: conversationID,
		RecipientID:    uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}
	
	return count, nil
}

// GetConversationParticipants gets all participants in a conversation
func (s *Service) GetConversationParticipants(ctx context.Context, conversationID uuid.UUID) ([]ParticipantResponse, error) {
	participants, err := s.store.GetConversationParticipants(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}
	
	return s.toParticipantResponses(participants), nil
}

// AddConversationParticipants adds multiple participants to a conversation
func (s *Service) AddConversationParticipants(ctx context.Context, conversationID uuid.UUID, userIDs []uuid.UUID) error {
	err := s.store.AddConversationParticipants(ctx, db.AddConversationParticipantsParams{
		ConversationID: conversationID,
		Column2:        userIDs,
	})
	if err != nil {
		return fmt.Errorf("failed to add participants: %w", err)
	}
	
	return nil
}

// LeaveConversation removes a user from a conversation
func (s *Service) LeaveConversation(ctx context.Context, conversationID, userID uuid.UUID) error {
	err := s.store.LeaveConversation(ctx, db.LeaveConversationParams{
		ConversationID: conversationID,
		UserID:         userID,
	})
	if err != nil {
		return fmt.Errorf("failed to leave conversation: %w", err)
	}
	
	return nil
}

// UpdateParticipantSettings updates participant preferences
func (s *Service) UpdateParticipantSettings(ctx context.Context, conversationID, userID uuid.UUID, req UpdateParticipantSettingsRequest) error {
	var notificationsEnabled sql.NullBool
	var customSettings pqtype.NullRawMessage
	
	if req.NotificationsEnabled != nil {
		notificationsEnabled = sql.NullBool{Bool: *req.NotificationsEnabled, Valid: true}
	}
	if req.CustomSettings != nil {
		customSettings = *req.CustomSettings
	}
	
	err := s.store.UpdateParticipantSettings(ctx, db.UpdateParticipantSettingsParams{
		ConversationID:       conversationID,
		UserID:               userID,
		NotificationsEnabled: notificationsEnabled,
		CustomSettings:       customSettings,
	})
	if err != nil {
		return fmt.Errorf("failed to update participant settings: %w", err)
	}
	
	return nil
}

// AddMessageReaction adds a reaction to a message
func (s *Service) AddMessageReaction(ctx context.Context, messageID uuid.UUID, req AddReactionRequest) error {
	// The SQLC function uses jsonb_set which requires specific parameters
	// For simplicity, we'll use a map to represent the reaction
	err := s.store.AddMessageReaction(ctx, db.AddMessageReactionParams{
		ID:      messageID,
		Column2: req.Emoji,
		ToJsonb: req.UserID.String(),
	})
	if err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}
	
	return nil
}

// RemoveMessageReaction removes a reaction from a message
func (s *Service) RemoveMessageReaction(ctx context.Context, messageID uuid.UUID, emoji string) error {
	// Note: The SQLC function signature needs the reactions field
	err := s.store.RemoveMessageReaction(ctx, db.RemoveMessageReactionParams{
		ID:        messageID,
		Reactions: pqtype.NullRawMessage{}, // This needs to be the emoji key to remove
	})
	if err != nil {
		return fmt.Errorf("failed to remove reaction: %w", err)
	}
	
	return nil
}

// Helper conversion functions

func (s *Service) toConversationResponse(c db.Conversation, unreadCount int64) *ConversationResponse {
	resp := &ConversationResponse{
		ID:          c.ID,
		SpaceID:     c.SpaceID,
		UnreadCount: unreadCount,
	}
	
	if c.Name.Valid {
		resp.Name = &c.Name.String
	}
	if c.Avatar.Valid {
		resp.Avatar = &c.Avatar.String
	}
	if c.Description.Valid {
		resp.Description = &c.Description.String
	}
	if c.ConversationType.Valid {
		resp.ConversationType = c.ConversationType.String
	}
	if c.LastMessageID.Valid {
		resp.LastMessageID = &c.LastMessageID.UUID
	}
	if c.LastMessageAt.Valid {
		resp.LastMessageAt = &c.LastMessageAt.Time
	}
	if c.IsActive.Valid {
		resp.IsActive = c.IsActive.Bool
	}
	if c.Settings.Valid {
		resp.Settings = &c.Settings
	}
	if c.CreatedAt.Valid {
		resp.CreatedAt = &c.CreatedAt.Time
	}
	if c.UpdatedAt.Valid {
		resp.UpdatedAt = &c.UpdatedAt.Time
	}
	
	return resp
}

func (s *Service) toConversationDetailResponse(c db.GetConversationByIDRow) *ConversationResponse {
	resp := &ConversationResponse{
		ID:          c.ID,
		SpaceID:     c.SpaceID,
		UnreadCount: c.UnreadCount,
	}
	
	if c.Name.Valid {
		resp.Name = &c.Name.String
	}
	if c.Avatar.Valid {
		resp.Avatar = &c.Avatar.String
	}
	if c.Description.Valid {
		resp.Description = &c.Description.String
	}
	if c.ConversationType.Valid {
		resp.ConversationType = c.ConversationType.String
	}
	if c.LastMessageID.Valid {
		resp.LastMessageID = &c.LastMessageID.UUID
	}
	if c.LastMessageAt.Valid {
		resp.LastMessageAt = &c.LastMessageAt.Time
	}
	if c.IsActive.Valid {
		resp.IsActive = c.IsActive.Bool
	}
	if c.Settings.Valid {
		resp.Settings = &c.Settings
	}
	if c.CreatedAt.Valid {
		resp.CreatedAt = &c.CreatedAt.Time
	}
	if c.UpdatedAt.Valid {
		resp.UpdatedAt = &c.UpdatedAt.Time
	}
	
	return resp
}

func (s *Service) toUserConversationResponses(conversations []db.GetUserConversationsRow) []ConversationDetailResponse {
	responses := make([]ConversationDetailResponse, len(conversations))
	for i, c := range conversations {
		responses[i] = s.toUserConversationResponse(c)
	}
	return responses
}

func (s *Service) toUserConversationResponse(c db.GetUserConversationsRow) ConversationDetailResponse {
	resp := ConversationDetailResponse{
		ID:          c.ID,
		SpaceID:     c.SpaceID,
		UnreadCount: c.UnreadCount,
	}
	
	if c.Name.Valid {
		resp.Name = &c.Name.String
	}
	if c.Avatar.Valid {
		resp.Avatar = &c.Avatar.String
	}
	if c.Description.Valid {
		resp.Description = &c.Description.String
	}
	if c.ConversationType.Valid {
		resp.ConversationType = c.ConversationType.String
	}
	if c.LastMessageID.Valid {
		resp.LastMessageID = &c.LastMessageID.UUID
	}
	if c.LastMessageAt.Valid {
		resp.LastMessageAt = &c.LastMessageAt.Time
	}
	if c.IsActive.Valid {
		resp.IsActive = c.IsActive.Bool
	}
	if c.Settings.Valid {
		resp.Settings = &c.Settings
	}
	if c.CreatedAt.Valid {
		resp.CreatedAt = &c.CreatedAt.Time
	}
	if c.UpdatedAt.Valid {
		resp.UpdatedAt = &c.UpdatedAt.Time
	}
	if c.UserRole.Valid {
		resp.UserRole = &c.UserRole.String
	}
	if c.NotificationsEnabled.Valid {
		resp.NotificationsEnabled = c.NotificationsEnabled.Bool
	}
	if c.LastMessageContent.Valid {
		resp.LastMessageContent = &c.LastMessageContent.String
	}
	if c.LastMessageTime.Valid {
		resp.LastMessageTime = &c.LastMessageTime.Time
	}
	if c.LastSenderUsername.Valid {
		resp.LastSenderUsername = &c.LastSenderUsername.String
	}
	if c.LastSenderFullName.Valid {
		resp.LastSenderFullName = &c.LastSenderFullName.String
	}
	
	return resp
}

func (s *Service) toBasicMessageResponse(m db.Message) *MessageResponse {
	resp := &MessageResponse{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		SenderID:       m.SenderID,
	}
	
	if m.RecipientID.Valid {
		resp.RecipientID = &m.RecipientID.UUID
	}
	if m.Content.Valid {
		resp.Content = m.Content.String
	}
	if m.Attachments.Valid {
		resp.Attachments = &m.Attachments
	}
	if m.MessageType.Valid {
		resp.MessageType = m.MessageType.String
	}
	if m.IsRead.Valid {
		resp.IsRead = m.IsRead.Bool
	}
	if m.ReadAt.Valid {
		resp.ReadAt = &m.ReadAt.Time
	}
	if m.Reactions.Valid {
		resp.Reactions = &m.Reactions
	}
	if m.ReplyToID.Valid {
		resp.ReplyToID = &m.ReplyToID.UUID
	}
	if m.Status.Valid {
		resp.Status = m.Status.String
	}
	if m.CreatedAt.Valid {
		resp.CreatedAt = &m.CreatedAt.Time
	}
	
	return resp
}

func (s *Service) toMessageResponse(m db.GetMessageByIDRow) *MessageResponse {
	resp := &MessageResponse{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		SenderID:       m.SenderID,
		SenderUsername: m.SenderUsername,
		SenderFullName: m.SenderFullName,
	}
	
	if m.RecipientID.Valid {
		resp.RecipientID = &m.RecipientID.UUID
	}
	if m.Content.Valid {
		resp.Content = m.Content.String
	}
	if m.Attachments.Valid {
		resp.Attachments = &m.Attachments
	}
	if m.MessageType.Valid {
		resp.MessageType = m.MessageType.String
	}
	if m.IsRead.Valid {
		resp.IsRead = m.IsRead.Bool
	}
	if m.ReadAt.Valid {
		resp.ReadAt = &m.ReadAt.Time
	}
	if m.Reactions.Valid {
		resp.Reactions = &m.Reactions
	}
	if m.ReplyToID.Valid {
		resp.ReplyToID = &m.ReplyToID.UUID
	}
	if m.Status.Valid {
		resp.Status = m.Status.String
	}
	if m.CreatedAt.Valid {
		resp.CreatedAt = &m.CreatedAt.Time
	}
	if m.SenderAvatar.Valid {
		resp.SenderAvatar = &m.SenderAvatar.String
	}
	
	return resp
}

func (s *Service) toMessageDetailResponses(messages []db.GetConversationMessagesRow) []MessageDetailResponse {
	responses := make([]MessageDetailResponse, len(messages))
	for i, m := range messages {
		responses[i] = s.toMessageDetailResponse(m)
	}
	return responses
}

func (s *Service) toMessageDetailResponse(m db.GetConversationMessagesRow) MessageDetailResponse {
	resp := MessageDetailResponse{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		SenderID:       m.SenderID,
		SenderUsername: m.SenderUsername,
		SenderFullName: m.SenderFullName,
	}
	
	if m.RecipientID.Valid {
		resp.RecipientID = &m.RecipientID.UUID
	}
	if m.Content.Valid {
		resp.Content = m.Content.String
	}
	if m.Attachments.Valid {
		resp.Attachments = &m.Attachments
	}
	if m.MessageType.Valid {
		resp.MessageType = m.MessageType.String
	}
	if m.IsRead.Valid {
		resp.IsRead = m.IsRead.Bool
	}
	if m.ReadAt.Valid {
		resp.ReadAt = &m.ReadAt.Time
	}
	if m.Reactions.Valid {
		resp.Reactions = &m.Reactions
	}
	if m.ReplyToID.Valid {
		resp.ReplyToID = &m.ReplyToID.UUID
	}
	if m.Status.Valid {
		resp.Status = m.Status.String
	}
	if m.CreatedAt.Valid {
		resp.CreatedAt = &m.CreatedAt.Time
	}
	if m.SenderAvatar.Valid {
		resp.SenderAvatar = &m.SenderAvatar.String
	}
	if m.ReplyContent.Valid {
		resp.ReplyContent = &m.ReplyContent.String
	}
	if m.ReplyUsername.Valid {
		resp.ReplyUsername = &m.ReplyUsername.String
	}
	
	return resp
}

func (s *Service) toParticipantResponses(participants []db.GetConversationParticipantsRow) []ParticipantResponse {
	responses := make([]ParticipantResponse, len(participants))
	for i, p := range participants {
		responses[i] = s.toParticipantResponse(p)
	}
	return responses
}

func (s *Service) toParticipantResponse(p db.GetConversationParticipantsRow) ParticipantResponse {
	resp := ParticipantResponse{
		ID:       p.ID,
		Username: p.Username,
		FullName: p.FullName,
	}
	
	if p.Avatar.Valid {
		resp.Avatar = &p.Avatar.String
	}
	if p.Verified.Valid {
		resp.Verified = p.Verified.Bool
	}
	if p.Role.Valid {
		resp.Role = p.Role.String
	}
	if p.JoinedAt.Valid {
		resp.JoinedAt = &p.JoinedAt.Time
	}
	if p.IsActive.Valid {
		resp.IsActive = p.IsActive.Bool
	}
	if p.NotificationsEnabled.Valid {
		resp.NotificationsEnabled = p.NotificationsEnabled.Bool
	}
	
	return resp
}
