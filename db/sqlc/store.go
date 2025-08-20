package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store cung cấp tất cả các hàm để thực thi truy vấn và transaction
type Store interface {
	Querier
}

// SQLStore triển khai interface Store
type SQLStore struct {
	*Queries
	db *pgxpool.Pool
}

// NewStore tạo một Store mới
func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}

	q := store.WithTx(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
