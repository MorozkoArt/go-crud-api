package repository

import (
    "context"
    "database/sql"
    "errors"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/MorozkoArt/go-crud-api/internal/auth"
    "github.com/MorozkoArt/go-crud-api/internal/models"
)

var (
    ErrUserNotFound = errors.New("user not found")
    ErrUserExists   = errors.New("user already exists")
)

type Repository struct {
    db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
    return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, u *models.User) error {
    var exists bool
    err := r.db.QueryRow(ctx, 
        "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", u.Email).
        Scan(&exists)
    if err != nil {
        return err
    }
    if exists {
        return ErrUserExists
    }

    hashedPassword, err := auth.HashPassword(u.Password)
    if err != nil {
        return err
    }

    _, err = r.db.Exec(ctx, 
        "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)",
        u.Name, u.Email, hashedPassword)
    return err
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    var u models.User
    err := r.db.QueryRow(ctx,
        "SELECT id, name, email, password FROM users WHERE email=$1", email).
        Scan(&u.Id, &u.Name, &u.Email, &u.Password)
    
    if err == sql.ErrNoRows {
        return nil, ErrUserNotFound
    }
    return &u, err
}

func (r *Repository) VerifyPassword(ctx context.Context, email, password string) (*models.User, error) {
    user, err := r.GetByEmail(ctx, email)
    if err != nil {
        return nil, err
    }

    if !auth.CheckPasswordHash(password, user.Password) {
        return nil, errors.New("invalid password")
    }

    user.Password = ""
    return user, nil
}

func (r *Repository) GetById(ctx context.Context, id int64) (*models.User, error) {
    var u models.User
    err := r.db.QueryRow(ctx,
        "SELECT id, name, email FROM users WHERE id=$1", id).
        Scan(&u.Id, &u.Name, &u.Email)
    
    if err == sql.ErrNoRows {
        return nil, ErrUserNotFound
    }
    return &u, err
}

func (r *Repository) GetAll(ctx context.Context) ([]models.User, error) {
    rows, err := r.db.Query(ctx, "SELECT id, name, email FROM users ORDER BY id")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var u models.User
        if err := rows.Scan(&u.Id, &u.Name, &u.Email); err != nil {
            return nil, err
        }
        users = append(users, u)
    }

    return users, nil
}

func (r *Repository) Update(ctx context.Context, u *models.User) error {
    result, err := r.db.Exec(ctx, 
        "UPDATE users SET name=$1, email=$2 WHERE id=$3",
        u.Name, u.Email, u.Id)
    if err != nil {
        return err
    }
    
    rowsAffected := result.RowsAffected()
    if rowsAffected == 0 {
        return ErrUserNotFound
    }
    return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
    result, err := r.db.Exec(ctx, "DELETE FROM users WHERE id=$1", id)
    if err != nil {
        return err
    }
    
    rowsAffected := result.RowsAffected()
    if rowsAffected == 0 {
        return ErrUserNotFound
    }
    return nil
}