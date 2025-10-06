package user

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, u *User) error {
	_, err := r.db.Exec(ctx, 
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3)",
		u.Name, u.Email, u.Password)
	return err
}