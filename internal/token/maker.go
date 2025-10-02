package token

import (
	"context"
	"time"
)

type Maker interface {
	CreateToken(userID int64, email string, roleId string, duration time.Duration) (string, error)
	CreateRefreshToken(userID int64, email string, roleId string, duration time.Duration) (string, error)
	VerifyToken(ctx context.Context, token string) (*Payload, error)
}
