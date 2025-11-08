package auth

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)


type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}


func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) < 32 {
		return nil, fmt.Errorf("invalid key size: must be at least 32 characters")
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}


func (maker *PasetoMaker) CreateToken(userID, username, spaceID string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, username, spaceID, duration)
	if err != nil {
		return "", nil, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, payload, err
}


func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}