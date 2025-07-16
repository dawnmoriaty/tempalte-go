package service

import (
	db "GIN/db/sqlc"
	"GIN/internal/repository"
	"context"
)

type UserService interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetAllUsers(ctx context.Context) ([]db.User, error)
}

type userServiceImpl struct {
	repo repository.UserRepository
}

// GetAllUsers implements UserService.
func (s *userServiceImpl) GetAllUsers(ctx context.Context) ([]db.User, error) {
	return s.repo.GetAllUsers(ctx)
}

// GetUserByEmail implements UserService.
func (s *userServiceImpl) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	// TODO: Hash password ở đây trước khi gọi repository
	return s.repo.CreateUser(ctx, arg)
}
