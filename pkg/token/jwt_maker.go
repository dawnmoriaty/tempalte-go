package token

import (
	"GIN/pkg/redis"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtMaker struct {
	secretKey         string
	redisTokenService redis.RedisTokenService
}

func (maker *JwtMaker) CreateToken(userID, username, email string, roles []string, tokenType string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, roles, duration)
	if err != nil {
		return "", payload, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", payload, err
	}
	tokenData := redis.TokenData{
		UserID:    userID,
		Username:  username,
		Email:     email,
		Roles:     roles,
		TokenType: tokenType,
	}
	err = maker.redisTokenService.StoreToken(payload.ID.String(), userID, tokenData, duration)
	if err != nil {
		return "", payload, fmt.Errorf("failed to store token in Redis: %w", err)
	}

	return tokenString, payload, nil
}

func (maker *JwtMaker) VerifyTokenWithRedis(tokenString string, tokenType string) (*Payload, *redis.TokenData, error) {
	payload, err := maker.VerifyToken(tokenString)
	if err != nil {
		return nil, nil, err
	}

	// 2. Kiểm tra token trong Redis
	tokenData, err := maker.redisTokenService.GetToken(payload.ID.String(), tokenType)
	if err != nil {
		return nil, nil, fmt.Errorf("token validation failed: %w", err)
	}

	// 3. Kiểm tra consistency giữa JWT payload và Redis data
	if payload.Username != tokenData.Username {
		return nil, nil, fmt.Errorf("token data inconsistency")
	}

	return payload, tokenData, nil
}

func (maker *JwtMaker) VerifyToken(tokenString string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInValidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(tokenString, &Payload{}, keyFunc)
	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInValidToken
	}

	return payload, nil
}

func (maker *JwtMaker) LogoutToken(tokenString, tokenType string) error {
	payload, err := maker.VerifyToken((tokenString))
	if err != nil {
		return fmt.Errorf("failed to verify token: %w", err)
	}
	remainingTTL := time.Until(payload.ExpiresAt.Time)
	if remainingTTL <= 0 {
		return nil // Token đã hết hạn
	}

	return maker.redisTokenService.BlacklistToken(payload.ID.String(), remainingTTL)
}

func (maker *JwtMaker) LogoutAllUserTokens(userID string) error {
	return maker.redisTokenService.RevokeAllUserTokens(userID)
}

func (maker *JwtMaker) RefreshTokenPair(oldRefreshToken string, userID, username, email string, roles []string, accessDuration, refreshDuration time.Duration) (string, string, error) {
	oldPayload, _, err := maker.VerifyTokenWithRedis(oldRefreshToken, "refresh")
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// 2. Create new token pair
	newAccessPayload, err := NewPayload(username, roles, accessDuration)
	if err != nil {
		return "", "", err
	}

	newRefreshPayload, err := NewPayload(username, roles, refreshDuration)
	if err != nil {
		return "", "", err
	}

	// 3. Generate new JWT strings
	newAccessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessPayload)
	newAccessToken, err := newAccessJWT.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", "", err
	}

	newRefreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, newRefreshPayload)
	newRefreshToken, err := newRefreshJWT.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", "", err
	}

	// 4. Update Redis with new token pair
	tokenData := redis.TokenData{
		UserID:   userID,
		Username: username,
		Email:    email,
		Roles:    roles,
	}

	err = maker.redisTokenService.RefreshToken(
		oldPayload.ID.String(),
		newAccessPayload.ID.String(),
		newRefreshPayload.ID.String(),
		userID,
		tokenData,
		accessDuration,
		refreshDuration,
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to refresh token in Redis: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

const minSecretKeyLength = 32

func NewJwtMaker(secretKey string, redisTokenService redis.RedisTokenService) (TokenMaker, error) {
	if len(secretKey) < minSecretKeyLength {
		return nil, fmt.Errorf("secret key length must be %d characters", minSecretKeyLength)
	}
	return &JwtMaker{secretKey: secretKey, redisTokenService: redisTokenService}, nil
}
