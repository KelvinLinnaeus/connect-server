package auth

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific user and duration
	CreateToken(userID, username string, spaceID string, duration time.Duration) (string, *Payload, error)

	// VerifyToken checks if the token is valid
	VerifyToken(token string) (*Payload, error)
}
