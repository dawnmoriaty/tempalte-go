package repository

import (
	db "GIN/db/sqlc"
	"context"
)

type RoleRepository interface {
	GetRoleByName(ctx context.Context, name string) (db.Role, error)
}

type roleRepositoryImpl struct {
	store db.Store
}

func NewRoleRepository(store db.Store) RoleRepository {
	return &roleRepositoryImpl{
		store: store,
	}
}

func (r *roleRepositoryImpl) GetRoleByName(ctx context.Context, name string) (db.Role, error) {
	return r.store.GetRoleByName(ctx, name)
}
