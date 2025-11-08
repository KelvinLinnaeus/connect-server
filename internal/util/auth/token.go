package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)


type HMACMaker struct {
	secretKey string
}


func NewHMACMaker(secret string) Maker {
	return &HMACMaker{secretKey: secret}
}


func (m *HMACMaker) CreateToken(userID, username, spaceID string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, username, spaceID, duration)
	if err != nil {
		return "", nil, err
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", nil, err
	}
	signature := m.sign(data)
	token := base64.URLEncoding.EncodeToString(data) + "." + base64.URLEncoding.EncodeToString(signature)
	return token, payload, nil
}


func (m *HMACMaker) VerifyToken(token string) (*Payload, error) {
	parts := splitToken(token)
	if len(parts) != 2 {
		return nil, errors.New("invalid token format")
	}
	dataB, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	sig, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	expected := m.sign(dataB)
	if !hmac.Equal(sig, expected) {
		return nil, errors.New("invalid token signature")
	}
	var p Payload
	if err := json.Unmarshal(dataB, &p); err != nil {
		return nil, err
	}
	if time.Now().After(p.ExpiredAt) {
		return nil, errors.New("token expired")
	}
	return &p, nil
}

func (m *HMACMaker) sign(data []byte) []byte {
	h := hmac.New(sha256.New, []byte(m.secretKey))
	h.Write(data)
	return h.Sum(nil)
}

func splitToken(tok string) []string {
	var parts []string
	for i := 0; i < len(tok); i++ {
		
		if tok[i] == '.' {
			parts = append(parts, tok[:i], tok[i+1:])
			return parts
		}
	}
	return []string{}
}
