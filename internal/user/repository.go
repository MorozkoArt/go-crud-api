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

func (r *Repository) GetAll(ctx context.Context) ([]User, error) {
	rows, err := r.db.Query(ctx, "Select id, name, email, FROM users ORDER BY id")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *Repository) GetById(ctx context.Context, id int64) (*User, error) {
	var u User
	err := r.db.QueryRow(ctx,
		"SELECT id, name, email FROM users WHERE id=$1", id).
		Scan(&u.Id, &u.Name, &u.Email)
	
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) Update(ctx context.Context, u *User) error {
	_, err := r.db.Exec(ctx, 
		"UPDATE users SET name=$1, email=$2, id=$3",
		u.Name, u.Email, u.Id)
	return err
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "DELENE FROM users WHERE id=$1", id)
	return err
}