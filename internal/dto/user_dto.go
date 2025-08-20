package dto

type LoginRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	UserName string `json:"user_name" binding:"required,min=4"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type JwtResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Type         string   `json:"type"`
	UserId       string   `json:"user_id"`
	UserName     string   `json:"user_name"`
	Email        string   `json:"email"`
	Role         []string `json:"role"`
}
