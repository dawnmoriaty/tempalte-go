package handler

import (
    db "GIN/db/sqlc"
    "GIN/internal/service"
    "net/http"

    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    userService service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
    return &UserHandler{userService: svc}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req db.CreateUserParams
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.userService.CreateUser(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, user)
}
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
    email := c.Param("email")
    user, err := h.userService.GetUserByEmail(c.Request.Context(), email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, user)
}
func (h *UserHandler) GetAllUsers(c *gin.Context) {
    users, err := h.userService.GetAllUsers(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, users)
}