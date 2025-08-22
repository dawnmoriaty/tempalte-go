package service

import (
	"GIN/configs"
	db "GIN/db/sqlc"
	"GIN/internal/dto"
	"GIN/internal/repository"
	"GIN/internal/utils"
	"GIN/pkg/token"
	"context"
	"errors"
	"time"
)

type UserService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (dto.JwtResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.JwtResponse, error)
	Logout(ctx context.Context, accessToken string) error
	LogoutAll(ctx context.Context, userID string) error
	RefreshToken(ctx context.Context, refreshToken string) (dto.JwtResponse, error)
	GetActiveSessionsCount(ctx context.Context, userID string) (int, error)
}

type userServiceImpl struct {
	repo       repository.UserRepository
	roleRepo   repository.RoleRepository
	tokenMaker token.TokenMaker
	config     *configs.Config
}

func (s *userServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (dto.JwtResponse, error) {
	existsEmail, err := s.repo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return dto.JwtResponse{}, err
	}
	if existsEmail {
		return dto.JwtResponse{}, errors.New("email already exists")
	}

	// Kiểm tra username đã tồn tại
	existsUserName, err := s.repo.CheckUserNameExists(ctx, req.UserName)
	if err != nil {
		return dto.JwtResponse{}, err
	}
	if existsUserName {
		return dto.JwtResponse{}, errors.New("username already exists")
	}

	// Hash password
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return dto.JwtResponse{}, err
	}
	arg := db.CreateUserParams{
		Username: req.UserName,
		Email:    req.Email,
		Password: hashPassword,
	}
	user, err := s.repo.CreateUser(ctx, arg)
	if err != nil {
		return dto.JwtResponse{}, err
	}

	// Gán role mặc định
	role, err := s.roleRepo.GetRoleByName(ctx, "ROLE_USER")
	if err != nil {
		return dto.JwtResponse{}, err
	}

	if err := s.repo.AddRoleToUser(ctx, db.AddRoleToUserParams{
		UserID: user.ID,
		RoleID: role.ID,
	}); err != nil {
		return dto.JwtResponse{}, err
	}

	// Lấy roles của user
	roles, err := s.repo.GetRolesForUser(ctx, user.ID)
	if err != nil {
		return dto.JwtResponse{}, err
	}
	accessDuration, err := time.ParseDuration(s.config.JWT.AccessTokenLife)
	if err != nil {
		return dto.JwtResponse{}, err
	}

	refreshDuration, err := time.ParseDuration(s.config.JWT.RefreshTokenLife)
	if err != nil {
		return dto.JwtResponse{}, err
	}

	accessToken, _, err := s.tokenMaker.CreateToken(
		user.ID.String(),
		user.Username,
		user.Email,
		roles,
		"access",
		accessDuration,
	)
	if err != nil {
		return dto.JwtResponse{}, err
	}

	refreshToken, _, err := s.tokenMaker.CreateToken(
		user.ID.String(),
		user.Username,
		user.Email,
		roles,
		"refresh",
		refreshDuration,
	)
	if err != nil {
		return dto.JwtResponse{}, err
	}

	return dto.JwtResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Type:         "Bearer",
		UserId:       user.ID.String(),
		UserName:     user.Username,
		Email:        user.Email,
		Role:         roles,
	}, nil
}

func (s *userServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (dto.JwtResponse, error) {
	user, err := s.repo.GetUserByUsername(ctx, req.UserName)
	if err != nil {
		return dto.JwtResponse{}, errors.New("invalid username or password")
	}

	// Kiểm tra password
	if err := utils.CheckPasswordHash(req.Password, user.Password); err != nil {
		return dto.JwtResponse{}, errors.New("invalid username or password")
	}

	// Lấy roles
	roles, err := s.repo.GetRolesForUser(ctx, user.ID)
	if err != nil {
		return dto.JwtResponse{}, err
	}

	if len(roles) == 0 {
		return dto.JwtResponse{}, errors.New("user has no roles assigned")
	}

	// Parse token durations
	accessDuration, err := time.ParseDuration(s.config.JWT.AccessTokenLife)
	if err != nil {
		return dto.JwtResponse{}, err
	}

	refreshDuration, err := time.ParseDuration(s.config.JWT.RefreshTokenLife)
	if err != nil {
		return dto.JwtResponse{}, err
	}

	// Tạo token pair với Redis
	accessToken, _, err := s.tokenMaker.CreateToken(
		user.ID.String(),
		user.Username,
		user.Email,
		roles,
		"access",
		accessDuration,
	)
	if err != nil {
		return dto.JwtResponse{}, err
	}

	refreshToken, _, err := s.tokenMaker.CreateToken(
		user.ID.String(),
		user.Username,
		user.Email,
		roles,
		"refresh",
		refreshDuration,
	)
	if err != nil {
		return dto.JwtResponse{}, err
	}

	return dto.JwtResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Type:         "Bearer",
		UserId:       user.ID.String(),
		UserName:     user.Username,
		Email:        user.Email,
		Role:         roles,
	}, nil
}

func (s *userServiceImpl) Logout(ctx context.Context, accessToken string) error {
	return s.tokenMaker.LogoutToken(accessToken, "access")
}

func (s *userServiceImpl) LogoutAll(ctx context.Context, userID string) error {
	return s.tokenMaker.LogoutAllUserTokens(userID)
}

func (s *userServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (dto.JwtResponse, error) {
	// Verify refresh token và lấy user info
	_, tokenData, err := s.tokenMaker.VerifyTokenWithRedis(refreshToken, "refresh")
	if err != nil {
		return dto.JwtResponse{}, errors.New("invalid refresh token")
	}

	// Parse durations
	accessDuration, err := time.ParseDuration(s.config.JWT.AccessTokenLife)
	if err != nil {
		return dto.JwtResponse{}, err
	}

	refreshDuration, err := time.ParseDuration(s.config.JWT.RefreshTokenLife)
	if err != nil {
		return dto.JwtResponse{}, err
	}

	// Tạo token pair mới
	newAccessToken, newRefreshToken, err := s.tokenMaker.RefreshTokenPair(
		refreshToken,
		tokenData.UserID,
		tokenData.Username,
		tokenData.Email,
		tokenData.Roles,
		accessDuration,
		refreshDuration,
	)
	if err != nil {
		return dto.JwtResponse{}, err
	}

	return dto.JwtResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		Type:         "Bearer",
		UserId:       tokenData.UserID,
		UserName:     tokenData.Username,
		Email:        tokenData.Email,
		Role:         tokenData.Roles,
	}, nil
}

func (s *userServiceImpl) GetActiveSessionsCount(ctx context.Context, userID string) (int, error) {
	return 0, nil
}

func NewUserService(role repository.UserRepository,
	roleRepo repository.RoleRepository,
	tokenMaker token.TokenMaker,
	config *configs.Config) UserService {
	return &userServiceImpl{
		repo:       role,
		roleRepo:   roleRepo,
		tokenMaker: tokenMaker,
		config:     config,
	}
}
