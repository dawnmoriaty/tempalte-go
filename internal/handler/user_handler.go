package handler

import (
	"GIN/internal/dto"
	"GIN/internal/middleware"
	"GIN/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler interface {
	CreateUser(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	LogoutAll(c *gin.Context)
	RefreshToken(c *gin.Context)
	GetProfile(c *gin.Context)
	GetActiveSessions(c *gin.Context)
}

type UserHandlerImpl struct {
	service service.UserService
}

func (h *UserHandlerImpl) CreateUser(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *UserHandlerImpl) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *UserHandlerImpl) Logout(c *gin.Context) {
	// Lấy access token từ header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorization header required"})
		return
	}

	// Extract token from "Bearer <token>"
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid authorization header format"})
		return
	}

	accessToken := authHeader[7:]
	err := h.service.Logout(c.Request.Context(), accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logout successful",
	})
}

func (h *UserHandlerImpl) LogoutAll(c *gin.Context) {
	// Lấy user data từ middleware context
	userData, exists := middleware.GetUserDataFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user data not found"})
		return
	}

	err := h.service.LogoutAll(c.Request.Context(), userData.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "all sessions logged out successfully",
	})
}

func (h *UserHandlerImpl) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *UserHandlerImpl) GetProfile(c *gin.Context) {
	// Lấy user data từ middleware context
	userData, exists := middleware.GetUserDataFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile information",
		"user": gin.H{
			"id":       userData.UserID,
			"username": userData.Username,
			"email":    userData.Email,
			"roles":    userData.Roles,
		},
	})
}

func (h *UserHandlerImpl) GetActiveSessions(c *gin.Context) {
	// Lấy user data từ middleware context
	userData, exists := middleware.GetUserDataFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user data not found"})
		return
	}

	count, err := h.service.GetActiveSessionsCount(c.Request.Context(), userData.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"active_sessions": count,
		"user_id":         userData.UserID,
	})
}

func NewUserHandler(s service.UserService) UserHandler {
	return &UserHandlerImpl{
		service: s,
	}
}
