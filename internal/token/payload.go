package token

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
	ErrTokenRevoked = errors.New("token has been revoked")
)

type Payload struct {
	ID             string    `json:"id"`
	UserId         int64     `json:"user_id"`
	Email          string    `json:"email"`
	RoleId         string    `json:"role_id"`
	SessionVersion int       `json:"session_version"`
	TokenType      string    `json:"token_type"`
	IssuedAt       time.Time `json:"iat"`
	ExpiredAt      time.Time `json:"exp"`
}

func NewPayload(userId int64, email, roleId string, sessionVersion int, duration time.Duration, tokenType string) (*Payload, error) {
	tokenId := fmt.Sprintf("token_%s_%d_%d", tokenType, userId, time.Now().Unix())

	payload := &Payload{
		ID:             tokenId,
		UserId:         userId,
		Email:          email,
		RoleId:         roleId,
		SessionVersion: sessionVersion,
		TokenType:      tokenType,
		IssuedAt:       time.Now(),
		ExpiredAt:      time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
