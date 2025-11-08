package auth

import "time"


type Maker interface {
	
	CreateToken(userID, username string, spaceID string, duration time.Duration) (string, *Payload, error)

	
	VerifyToken(token string) (*Payload, error)
}
