package redis

import (
	"GIN/configs"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

type TokenData struct {
	UserID    string   `json:"user_id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	TokenType string   `json:"token_type"`
}

type RedisTokenServiceImpl struct {
	client *redis.Client
	ctx    context.Context
}

func (r *RedisTokenServiceImpl) StoreToken(tokenID, userID string, tokenData TokenData, ttl time.Duration) error {
	// lưu tokne với ttl
	tokenKey := fmt.Sprintf("token:%s:%s", tokenData.TokenType, tokenID)
	tokenDataJson, err := json.Marshal(tokenData)
	if err != nil {
		return fmt.Errorf("token data json marshal fail")
	}
	err = r.client.Set(r.ctx, tokenKey, tokenDataJson, ttl).Err()
	if err != nil {
		return fmt.Errorf("token data set fail : %w", err)
	}
	// thêm token vào session list của user
	userSessionKey := fmt.Sprintf("user:%s:session", userID)
	err = r.client.SAdd(r.ctx, userSessionKey, tokenKey).Err()
	if err != nil {
		return fmt.Errorf("failed to add token to user sessions: %w", err)
	}
	r.client.Expire(r.ctx, userSessionKey, ttl+time.Hour)
	return nil
}

func (r *RedisTokenServiceImpl) GetToken(tokenID, tokenType string) (*TokenData, error) {
	tokenKey := fmt.Sprintf("token:%s:%s", tokenType, tokenID)
	if r.IsTokenBlacklisted(tokenID) {
		return nil, fmt.Errorf("token id %s is blacklisted", tokenID)
	}
	tokenDataJson, err := r.client.Get(r.ctx, tokenKey).Result()
	if err != nil {
		return nil, fmt.Errorf("token data get fail : %w", err)
	}
	var tokenData TokenData
	err = json.Unmarshal([]byte(tokenDataJson), &tokenData)
	if err != nil {
		return nil, fmt.Errorf("token data json unmarshal fail : %w", err)
	}
	return &tokenData, nil
}

func (r *RedisTokenServiceImpl) BlacklistToken(tokenID string, remainingTTL time.Duration) error {
	blackListKey := fmt.Sprintf("blacklist:token:%s", tokenID)
	return r.client.Set(r.ctx, blackListKey, time.Now().Unix(), remainingTTL).Err()
}

func (r *RedisTokenServiceImpl) IsTokenBlacklisted(tokenID string) bool {
	blackListKey := fmt.Sprintf("blacklist:token:%s", tokenID)
	exists, _ := r.client.Exists(r.ctx, blackListKey).Result()
	return exists > 0
}

func (r *RedisTokenServiceImpl) RevokeAllUserTokens(userID string) error {
	userSessionKey := fmt.Sprintf("user:%s:session", userID)
	tokenIDs, err := r.client.SMembers(r.ctx, userSessionKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get user sessions: %w", err)
	}
	pipe := r.client.Pipeline()
	for _, tokenID := range tokenIDs {
		accessTokenKey := fmt.Sprintf("access_token:%s", tokenID)
		refreshTokenKey := fmt.Sprintf("refresh_token:%s", tokenID)

		// Blacklist cả access và refresh token
		pipe.Set(r.ctx, fmt.Sprintf("blacklist:%s", tokenID), time.Now().Unix(), time.Hour*24)

		// Xóa token khỏi Redis
		pipe.Del(r.ctx, accessTokenKey, refreshTokenKey)
	}
	pipe.Del(r.ctx, userSessionKey)
	return err
}

func (r *RedisTokenServiceImpl) RefreshToken(oldRefreshTokenID, newAccessTokenID, newRefreshTokenID, userID string, tokenData TokenData, accessTTL, refreshTTL time.Duration) error {
	pipe := r.client.Pipeline()
	pipe.Set(r.ctx, fmt.Sprintf("blacklist:%s", oldRefreshTokenID), time.Now().Unix(), refreshTTL)
	oldRefreshKey := fmt.Sprintf("refresh_token:%s", oldRefreshTokenID)
	pipe.Del(r.ctx, oldRefreshKey, oldRefreshTokenID)
	accessTokenData := tokenData
	accessTokenData.TokenType = AccessTokenType
	accessDataJson, _ := json.Marshal(accessTokenData)
	pipe.Set(r.ctx, fmt.Sprintf("access_token:%s", newAccessTokenID), accessDataJson, accessTTL)

	refreshTokenData := tokenData
	refreshTokenData.TokenType = RefreshTokenType
	refreshDataJson, _ := json.Marshal(refreshTokenData)
	pipe.Set(r.ctx, fmt.Sprintf("refresh_token:%s", newRefreshTokenID), refreshDataJson, refreshTTL)

	userSessionKey := fmt.Sprintf("user_sessions:%s", userID)
	pipe.SRem(r.ctx, userSessionKey, oldRefreshTokenID)
	pipe.SAdd(r.ctx, userSessionKey, newAccessTokenID, newRefreshTokenID)
	pipe.Expire(r.ctx, userSessionKey, refreshTTL+time.Hour)

	_, err := pipe.Exec(r.ctx)
	return err
}

func (r *RedisTokenServiceImpl) GetActiveUserSessions(userID string) (int, error) {
	userSessionKey := fmt.Sprintf("user:%s:session", userID)
	count, err := r.client.SCard(r.ctx, userSessionKey).Result()
	return int(count), err
}

func (r *RedisTokenServiceImpl) CleanupExpiredTokens() error {
	return nil
}

func NewRedisTokenService() RedisTokenService {
	config := configs.GetConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	return &RedisTokenServiceImpl{
		client: client,
		ctx:    context.Background(),
	}
}

type RedisTokenService interface {
	StoreToken(tokenID, userID string, tokenData TokenData, ttl time.Duration) error
	GetToken(tokenID, tokenType string) (*TokenData, error)
	BlacklistToken(tokenID string, remainingTTL time.Duration) error
	IsTokenBlacklisted(tokenID string) bool
	RevokeAllUserTokens(userID string) error
	RefreshToken(oldRefreshTokenID, newAccessTokenID, newRefreshTokenID, userID string, tokenData TokenData, accessTTL, refreshTTL time.Duration) error
	GetActiveUserSessions(userID string) (int, error)
	CleanupExpiredTokens() error
}
