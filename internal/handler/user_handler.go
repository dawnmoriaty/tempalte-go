package handler

import (
	db "GIN/db/sqlc"
	"GIN/internal/service"
	"GIN/pkg/response"
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
        response.SendError(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
        return
    }

    user, err := h.userService.CreateUser(c.Request.Context(), req)
    if err != nil {
        response.SendError(c, http.StatusInternalServerError, "Failed to create user: "+err.Error())
        return
    }

    response.SendSuccess(c, "User created successfully", user)
}
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
    email := c.Param("email")
    user, err := h.userService.GetUserByEmail(c.Request.Context(), email)
    if err != nil {
        response.SendError(c, http.StatusInternalServerError, "Failed to retrieve user: "+err.Error())
        return
    }

    response.SendSuccess(c, "User retrieved successfully", user)
}
func (h *UserHandler) GetAllUsers(c *gin.Context) {
    users, err := h.userService.GetAllUsers(c.Request.Context())
    if err != nil {
        response.SendError(c, http.StatusInternalServerError, "Failed to retrieve users: "+err.Error())
        return
    }

    response.SendSuccess(c, "Users retrieved successfully", users)
}