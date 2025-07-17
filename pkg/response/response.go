package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func SendSuccess[T any](c *gin.Context, message string, data T) {
	c.JSON(http.StatusOK, APIResponse[T]{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SendError(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, APIResponse[any]{
		Success: false,
		Message: message,
		Data:    nil,
	})
}