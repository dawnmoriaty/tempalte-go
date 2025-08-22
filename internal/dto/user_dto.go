package dto

type RegisterRequest struct {
	UserName string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type JwtResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Type         string   `json:"token_type"`
	UserId       string   `json:"user_id"`
	UserName     string   `json:"username"`
	Email        string   `json:"email"`
	Role         []string `json:"roles"`
	ExpiresIn    int64    `json:"expires_in,omitempty"` // Seconds until access token expires
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	LogoutAll bool `json:"logout_all,omitempty"` // If true, logout from all devices
}

type SessionInfo struct {
	TokenID   string `json:"token_id"`
	UserAgent string `json:"user_agent,omitempty"`
	IPAddress string `json:"ip_address,omitempty"`
	LoginTime int64  `json:"login_time"`
	LastUsed  int64  `json:"last_used,omitempty"`
	Device    string `json:"device,omitempty"`
}

type ActiveSessionsResponse struct {
	TotalSessions    int           `json:"total_sessions"`
	CurrentSessionID string        `json:"current_session_id"`
	Sessions         []SessionInfo `json:"sessions"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
