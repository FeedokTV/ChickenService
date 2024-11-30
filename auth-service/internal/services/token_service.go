package services

import (
	logger "auth-service/internal"
	"auth-service/internal/auth"
	"auth-service/internal/domain"
	"auth-service/internal/utils"
	"context"
	"time"

	"go.uber.org/zap"
)

type TokenService struct {
	tokenRepo domain.TokenRepository
}

func NewTokenService(tokenRepo domain.TokenRepository) *TokenService {
	return &TokenService{tokenRepo: tokenRepo}
}

func (s *TokenService) CreateToken(ctx context.Context, userID int, fingerprint string) (*domain.Token, *utils.APIError) {
	// Generate new JWT token
	tokenString, err := auth.GenerateToken(userID)
	if err != nil {
		logger.Error("Failed to generate token",
			zap.Int("User ID", userID),
			zap.Error(err))
		return nil, utils.NewAPIError(500, "Failed to generate token", "")
	}

	newToken := &domain.Token{
		UserID:    userID,
		Token:     tokenString,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(12 * time.Hour),
	}

	// Save token in repo
	apiErr := s.tokenRepo.SaveToken(ctx, userID, fingerprint, newToken)
	if apiErr != nil {
		logger.Error("Cannot save token in Redis.",
			zap.String("error", apiErr.Message),
			zap.String("details", apiErr.Details),
			zap.Int("User ID", userID))
		return nil, apiErr
	}

	return newToken, nil
}

func (s *TokenService) ValidateToken(ctx context.Context, token string, fingerprint string) (*auth.Claims, *utils.APIError) {

	claims, err := auth.ValidateToken(token)
	if err != nil {
		return nil, utils.NewAPIError(403, "Invalid token", "")
	}

	// Is expired token
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, utils.NewAPIError(403, "Expired token", "")
	}

	sessionToken, apiErr := s.tokenRepo.GetToken(context.Background(), claims.UserID, fingerprint)

	// Token expired by TTL or not found in redis
	if apiErr != nil {
		return nil, utils.NewAPIError(403, "Invalid or expired token", "")
	}

	// Token for current fingerprint session not the same as passed token
	if sessionToken != token {
		return nil, utils.NewAPIError(403, "Invalid or expired token", "")
	}

	return claims, nil
}
