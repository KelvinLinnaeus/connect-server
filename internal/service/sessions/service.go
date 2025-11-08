package sessions

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	"github.com/google/uuid"
)

// Service handles session-related business logic
type Service struct {
	store db.Store
}

// NewService creates a new sessions service
func NewService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}

// GetSession retrieves a session by ID
func (s *Service) GetSession(ctx context.Context, sessionID uuid.UUID) (*SessionResponse, error) {
	session, err := s.store.GetSession(ctx, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return s.toSessionResponse(session), nil
}

// CreateSession creates a new session
func (s *Service) CreateSession(ctx context.Context, req CreateSessionRequest) (*SessionResponse, error) {
	sessionID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	var ipAddress sql.NullString
	if req.IPAddress != nil {
		ipAddress = sql.NullString{String: *req.IPAddress, Valid: true}
	}

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           sessionID,
		UserID:       req.UserID,
		Username:     req.Username,
		RefreshToken: req.RefreshToken,
		UserAgent:    req.UserAgent,
		IpAddress:    ipAddress,
		IsBlocked:    false,
		SpaceID:      req.SpaceID,
		LastActivity: sql.NullTime{Time: time.Now(), Valid: true},
		ExpiresAt:    req.ExpiresAt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return s.toSessionResponse(session), nil
}

// toSessionResponse converts a database UserSession to SessionResponse
func (s *Service) toSessionResponse(session db.UserSession) *SessionResponse {
	resp := &SessionResponse{
		ID:        session.ID,
		UserID:    session.UserID,
		SpaceID:   session.SpaceID,
		Username:  session.Username,
		UserAgent: session.UserAgent,
		IsBlocked: session.IsBlocked,
		ExpiresAt: session.ExpiresAt,
	}

	if session.IpAddress.Valid {
		resp.IPAddress = &session.IpAddress.String
	}

	if session.LastActivity.Valid {
		resp.LastActivity = &session.LastActivity.Time
	}

	if session.CreatedAt.Valid {
		resp.CreatedAt = &session.CreatedAt.Time
	}

	return resp
}
