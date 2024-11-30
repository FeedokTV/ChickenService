package repositories

import (
	"auth-service/internal/domain"
	"auth-service/internal/utils"
	"context"
	"fmt"
	"time"

	logger "auth-service/internal"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisTokenRepo struct {
	client *redis.Client
}

func NewRedisTokenRepo(client *redis.Client) *RedisTokenRepo {
	return &RedisTokenRepo{client: client}
}

func (repo *RedisTokenRepo) SaveToken(ctx context.Context, userID int, fingerprintHash string, token *domain.Token) *utils.APIError {
	key := fmt.Sprintf("%d:%s", userID, fingerprintHash)
	err := repo.client.Set(ctx, key, token.Token, time.Until(token.ExpiresAt)).Err()
	if err != nil {
		logger.Error("Cannot save token",
			zap.Error(err))
		return utils.NewAPIError(500, "Failed to save token", err.Error())
	}

	return nil
}

func (repo *RedisTokenRepo) GetToken(ctx context.Context, userID int, fingerprintHash string) (string, *utils.APIError) {
	key := fmt.Sprintf("%d:%s", userID, fingerprintHash)
	token, err := repo.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", utils.NewAPIError(404, "Token not found or expired", "")
		}
		logger.Error("Cannot get token",
			zap.Error(err))
		return "", utils.NewAPIError(500, "Failed to get token", err.Error())
	}

	return token, nil
}

func (repo *RedisTokenRepo) IsTokenExists(ctx context.Context, token string) bool {

	_, err := repo.client.Get(ctx, token).Result()
	if err != nil {
		if err == redis.Nil {
			return false
		}
		logger.Error("Cannot get token",
			zap.Error(err))
		return false
	}

	return true
}
