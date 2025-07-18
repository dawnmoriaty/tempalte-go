package service

import (
	db "GIN/db/sqlc"
	"GIN/internal/repository"
	"context"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, params CreateUserParams) (db.User, error)
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

func (s *userServiceImpl) CreateUser(ctx context.Context, params CreateUserParams) (db.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if  err != nil {
		return db.User{}, err
	}
	repoParams := db.CreateUserParams{
		Email:    params.Email,
		Name:     params.Name,
		Password: string(hashedPassword),
	}
	return s.repo.CreateUser(ctx, repoParams)
}

type CreateUserParams struct {
	Email   string
	Name    string
	Password string
}
