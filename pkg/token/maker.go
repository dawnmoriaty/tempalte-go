package token

import (
	"GIN/pkg/redis"
	"time"
)

type TokenMaker interface {
	CreateToken(userID, username, email string, roles []string, tokenType string, duration time.Duration) (string, *Payload, error)
	VerifyTokenWithRedis(tokenString string, tokenType string) (*Payload, *redis.TokenData, error)
	VerifyToken(tokenString string) (*Payload, error)
	LogoutToken(tokenString, tokenType string) error
	LogoutAllUserTokens(userID string) error
	RefreshTokenPair(oldRefreshToken string, userID, username, email string, roles []string, accessDuration, refreshDuration time.Duration) (string, string, error)
}
