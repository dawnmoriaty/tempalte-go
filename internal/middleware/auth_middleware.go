package middleware

import (
	"GIN/pkg/redis"
	"GIN/pkg/token"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
	UserDataKey             = "user_data" // Key để lưu user data từ Redis
)

func AuthMiddlewareWithRedis(tokenMaker token.TokenMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy header
		authorizationHeader := c.Request.Header.Get(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Tách header
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Kiểm tra type
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			err := errors.New("unsupported authorization type")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Xác thực token với Redis
		accessToken := fields[1]
		payload, tokenData, err := tokenMaker.VerifyTokenWithRedis(accessToken, "access")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid or expired token",
				"details": err.Error(),
			})
			return
		}

		// Lưu thông tin vào context
		c.Set(AuthorizationPayloadKey, payload)
		c.Set(UserDataKey, tokenData)
		c.Next()
	}
}

// Helper function để lấy user data từ context
func GetUserDataFromContext(c *gin.Context) (*redis.TokenData, bool) {
	userData, exists := c.Get(UserDataKey)
	if !exists {
		return nil, false
	}

	tokenData, ok := userData.(*redis.TokenData)
	return tokenData, ok
}

func GetPayloadFromContext(c *gin.Context) (*token.Payload, bool) {
	payload, exists := c.Get(AuthorizationPayloadKey)
	if !exists {
		return nil, false
	}

	tokenPayload, ok := payload.(*token.Payload)
	return tokenPayload, ok
}
