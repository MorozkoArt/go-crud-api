package db

import (
	"context"
	"fmt"
	"log"
	
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/MorozkoArt/go-crud-api/internal/config"
)

func NewPostgreDB(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("Connection to PostgreSQL established")
	return pool, nil
}
