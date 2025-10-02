package token

import (
	"context"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type Paseto struct {
	paseto         *paseto.V2
	symmetric      []byte
	userRepository UserRepository
}

type UserRepository interface {
	GetSessionVersion(ctx context.Context, userID int64) (int, error)
}

func NewPaseto(symmetricKey string, userRepository UserRepository) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &Paseto{
		paseto:         paseto.NewV2(),
		symmetric:      []byte(symmetricKey),
		userRepository: userRepository,
	}

	return maker, nil
}

func (maker *Paseto) CreateToken(userID int64, email string, roleId string, duration time.Duration) (string, error) {
	sv := 1
	if maker.userRepository != nil {
		if v, err := maker.userRepository.GetSessionVersion(context.Background(), userID); err == nil {
			sv = v
		}
	}
	payload, err := NewPayload(userID, email, roleId, sv, duration, "access")
	if err != nil {
		return "", fmt.Errorf("failed to create payload: %w", err)
	}

	token, err := maker.paseto.Encrypt(maker.symmetric, payload, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt refresh token: %w", err)
	}

	return token, nil
}

func (maker *Paseto) CreateRefreshToken(userID int64, email string, roleId string, duration time.Duration) (string, error) {
	sv := 1
	if maker.userRepository != nil {
		if v, err := maker.userRepository.GetSessionVersion(context.Background(), userID); err == nil {
			sv = v
		}
	}
	payload, err := NewPayload(userID, email, roleId, sv, duration, "refresh")
	if err != nil {
		return "", fmt.Errorf("failed to create payload: %w", err)
	}

	token, err := maker.paseto.Encrypt(maker.symmetric, payload, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt refresh token: %w", err)
	}

	return token, nil
}

func (maker *Paseto) VerifyToken(ctx context.Context, token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetric, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	// check if token is expired
	if time.Now().After(payload.ExpiredAt) {
		return nil, ErrExpiredToken
	}

	if maker.userRepository != nil {
		currentSessionVersion, err := maker.userRepository.GetSessionVersion(ctx, payload.UserId)
		if err == nil && currentSessionVersion != payload.SessionVersion {
			return nil, ErrTokenRevoked // Session invalidated
		}
	}

	return payload, nil
}
