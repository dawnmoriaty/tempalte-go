package repository

import (
	db "GIN/db/sqlc"
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetAllUsers(ctx context.Context) ([]db.User, error)
}

type userRepositoryImpl struct {
	store db.Store
}

// GetAllUsers implements UserRepository.
func (r *userRepositoryImpl) GetAllUsers(ctx context.Context) ([]db.User, error) {
	return r.store.GetAllUsers(ctx)
}

// GetUserByEmail implements UserRepository.
func (r *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.store.GetUserByEmail(ctx, email)
}

func NewUserRepository(store db.Store) UserRepository {
	return &userRepositoryImpl{store: store}
}

func (r *userRepositoryImpl) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return r.store.CreateUser(ctx, arg)
}
