package database

import (
	"GIN/configs"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(cfg *configs.DatabaseConfig) *pgxpool.Pool {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	pool, err := pgxpool.New( context.Background(),dsn )
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v", err))
	}
	if err := pool.Ping(context.Background()); err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}

	log.Println("✅ Database connection successful!")
	return pool
}