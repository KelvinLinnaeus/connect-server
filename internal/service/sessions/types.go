package sessions

import (
	"time"

	"github.com/google/uuid"
)


type SessionResponse struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	SpaceID      uuid.UUID  `json:"space_id"`
	Username     string     `json:"username"`
	UserAgent    string     `json:"user_agent"`
	IPAddress    *string    `json:"ip_address,omitempty"`
	IsBlocked    bool       `json:"is_blocked"`
	LastActivity *time.Time `json:"last_activity,omitempty"`
	ExpiresAt    time.Time  `json:"expires_at"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
}


type CreateSessionRequest struct {
	UserID       uuid.UUID `json:"user_id" binding:"required"`
	Username     string    `json:"username" binding:"required"`
	RefreshToken string    `json:"refresh_token" binding:"required"`
	UserAgent    string    `json:"user_agent"`
	IPAddress    *string   `json:"ip_address"`
	SpaceID      uuid.UUID `json:"space_id" binding:"required"`
	ExpiresAt    time.Time `json:"expires_at" binding:"required"`
}


type ListSessionsRequest struct {
	UserID uuid.UUID `json:"user_id"`
}
