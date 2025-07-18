package dto

import (
	db "GIN/db/sqlc"
	"time"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100" log:"-"`
}

type UpdateUserRequest struct {
	Name *string `json:"name,omitempty" binding:"omitempty,min=3,max=50"`
}

type UserResponse struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	RoleID    int32   `json:"role_id,omitempty"`
	AvatarUrl string  `json:"avatar_url,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

func MapUserToResponse(user db.User) UserResponse {
	var deletedAt *string
	if user.DeletedAt.Valid {
		deletedAtStr := user.DeletedAt.Time.Format(time.RFC3339)
		deletedAt = &deletedAtStr
	}

	return UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		AvatarUrl: user.AvatarUrl.String,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		DeletedAt: deletedAt,
		RoleID:    user.RoleID,
	}
}