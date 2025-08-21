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
}

type userServiceImpl struct {
	repo       repository.UserRepository
	roleRepo   repository.RoleRepository
	tokenMaker token.TokenMaker
	config     *configs.Config
}

func (s *userServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (dto.JwtResponse, error) {
	//TODO implement me
	exitsEmail, err := s.repo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return dto.JwtResponse{}, err
	}
	if exitsEmail {
		return dto.JwtResponse{}, errors.New("email already exists")
	}
	existsUserName, err := s.repo.CheckUserNameExists(ctx, req.UserName)
	if err != nil {
		return dto.JwtResponse{}, err
	}
	if existsUserName {
		return dto.JwtResponse{}, errors.New("username already exists")
	}
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
	role, err := s.roleRepo.GetRoleByName(ctx, "ROLE_USER")
	if err != nil {
		return dto.JwtResponse{}, err
	}
	if err := s.repo.AddRoleToUser(
		ctx, db.AddRoleToUserParams{
			UserID: user.ID,
			RoleID: role.ID,
		}); err != nil {
		return dto.JwtResponse{}, err
	}
	roles, err := s.repo.GetRolesForUser(ctx, user.ID)
	if err != nil {
		return dto.JwtResponse{}, err
	}
	accessionDuration, err := time.ParseDuration(s.config.JWT.AccessTokenLife)
	if err != nil {
		return dto.JwtResponse{}, err
	}
	accessToken, _, err := s.tokenMaker.CreateToken(user.Username, roles, accessionDuration)
	if err != nil {
		return dto.JwtResponse{}, err
	}
	refreshTokenDuration, err := time.ParseDuration(s.config.JWT.RefreshTokenLife)
	if err != nil {
		return dto.JwtResponse{}, err
	}
	refreshToken, _, err := s.tokenMaker.CreateToken(user.Username, roles, refreshTokenDuration)
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
	if err := utils.CheckPasswordHash(req.Password, user.Password); err != nil {
		return dto.JwtResponse{}, errors.New("invalid username or password")
	}
	roles, err := s.repo.GetRolesForUser(ctx, user.ID)
	if err != nil {
		return dto.JwtResponse{}, err
	}
	if len(roles) == 0 {
		return dto.JwtResponse{}, errors.New("user has no roles assigned")
	}
	accessTokenDuration, err := time.ParseDuration(s.config.JWT.AccessTokenLife)
	if err != nil {
		return dto.JwtResponse{}, err
	}
	accessToken, _, err := s.tokenMaker.CreateToken(user.Username, roles, accessTokenDuration)
	if err != nil {
		return dto.JwtResponse{}, err
	}
	refreshTokenDuration, err := time.ParseDuration(s.config.JWT.RefreshTokenLife)
	if err != nil {
		return dto.JwtResponse{}, err
	}
	refreshToken, _, err := s.tokenMaker.CreateToken(user.Username, roles, refreshTokenDuration)
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
