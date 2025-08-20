package handler

import (
	"GIN/internal/dto"
	"GIN/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler interface {
	CreateUser(c *gin.Context)
	Login(c *gin.Context)
}

type UserHandlerImpl struct {
	service service.UserService
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

func (h *UserHandlerImpl) CreateUser(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	res, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func NewUserHandler(s service.UserService) UserHandler {
	return &UserHandlerImpl{
		service: s,
	}
}
