package handler

import (
	"GIN/internal/dto"
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
    var req dto.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.SendError(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
        return
    }
    params:= service.CreateUserParams{
        Email:    req.Email,
        Name:     req.Name,
        Password: req.Password,
    }
    user, err := h.userService.CreateUser(c.Request.Context(), params)
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
    responseData := dto.MapUserToResponse(user)
    response.SendSuccess(c, "User retrieved successfully", responseData)
}
func (h *UserHandler) GetAllUsers(c *gin.Context) {
    users, err := h.userService.GetAllUsers(c.Request.Context())
    if err != nil {
        response.SendError(c, http.StatusInternalServerError, "Failed to retrieve users: "+err.Error())
        return
    }
    var userResponses []dto.UserResponse
    for _, user := range users {
        userResponses = append(userResponses, dto.MapUserToResponse(user))
    }

    response.SendSuccess(c, "Users retrieved successfully", userResponses)

}