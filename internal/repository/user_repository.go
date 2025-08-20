package repository

import (
	db "GIN/db/sqlc"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	CheckUserNameExists(ctx context.Context, username string) (bool, error)
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	AddRoleToUser(ctx context.Context, arg db.AddRoleToUserParams) error
	GetRolesForUser(ctx context.Context, userID uuid.UUID) ([]string, error)
	GetUserByUsername(ctx context.Context, username string) (db.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (db.User, error)
}

type userRepositoryImpl struct {
	store db.Store
}

func (u *userRepositoryImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (db.User, error) {
	return u.store.GetUserById(ctx, userID)
}

func (u *userRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (db.User, error) {
	return u.store.GetUserByUsername(ctx, username)
}

func (u *userRepositoryImpl) CheckUserNameExists(ctx context.Context, username string) (bool, error) {
	return u.store.CheckUserNameExists(ctx, username)
}

// AddRoleToUser implements UserRepository.
// Subtle: this method shadows the method (*Queries).AddRoleToUser of userRepositoryImpl.Queries.
func (u *userRepositoryImpl) AddRoleToUser(ctx context.Context, arg db.AddRoleToUserParams) error {
	return u.store.AddRoleToUser(ctx, arg)
}

// CheckEmailExists implements UserRepository.
// Subtle: this method shadows the method (*Queries).CheckEmailExists of userRepositoryImpl.Queries.
func (u *userRepositoryImpl) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	return u.store.CheckEmailExists(ctx, email)
}

// CreateUser implements UserRepository.
// Subtle: this method shadows the method (*Queries).CreateUser of userRepositoryImpl.Queries.
func (u *userRepositoryImpl) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return u.store.CreateUser(ctx, arg)
}

func (u *userRepositoryImpl) GetRolesForUser(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return u.store.GetRolesForUser(ctx, userID)
}

// GetUserByEmail implements UserRepository.
// Subtle: this method shadows the method (*Queries).GetUserByEmail of userRepositoryImpl.Queries.
func (u *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return u.store.GetUserByEmail(ctx, email)
}

func NewUserRepository(store db.Store) UserRepository {
	return &userRepositoryImpl{
		store: store,
	}
}
