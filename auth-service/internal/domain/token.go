package domain

import (
	"auth-service/internal/utils"
	"context"
	"time"
)

type (
	TokenRepository interface {
		SaveToken(ctx context.Context, userID int, fingerprintHash string, token *Token) *utils.APIError
		GetToken(ctx context.Context, userID int, fingerprintHash string) (string, *utils.APIError)
		IsTokenExists(ctx context.Context, token string) bool
	}

	Token struct {
		UserID    int       `json:"user_id"`
		Token     string    `json:"token"`
		IssuedAt  time.Time `json:"-"`
		ExpiresAt time.Time `json:"-"`
	}
)
