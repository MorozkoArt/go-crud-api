package repository

import (
    "context"
    "database/sql"
    "errors"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/MorozkoArt/go-crud-api/internal/models"
    "github.com/MorozkoArt/go-crud-api/internal/utils"
)

var (
    ErrUserNotFound = errors.New("user not found")
    ErrUserExists   = errors.New("user already exists")
)

type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    GetByID(ctx context.Context, id int64) (*models.User, error)
    GetAll(ctx context.Context) ([]models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id int64) error
}

type userRepository struct {
    db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *models.User) error {
    log.Printf("Creating user with email: %s", u.Email)
    
    var exists bool
    err := r.db.QueryRow(ctx, 
        "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", u.Email).
        Scan(&exists)
    if err != nil {
        log.Printf("Error checking user existence: %v", err)
        return err
    }
    if exists {
        log.Printf("User with email %s already exists", u.Email)
        return ErrUserExists
    }

    hashedPassword, err := utils.HashPassword(u.Password)
    if err != nil {
        log.Printf("Error hashing password: %v", err)
        return err
    }

    _, err = r.db.Exec(ctx, 
        "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)",
        u.Name, u.Email, hashedPassword)
        
    if err != nil {
        log.Printf("Error creating user: %v", err)
    } else {
        log.Printf("User created successfully: %s", u.Email)
    }
    
    return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    log.Printf("Fetching user by email: %s", email)
    
    var u models.User
    err := r.db.QueryRow(ctx,
        "SELECT id, name, email, password FROM users WHERE email=$1", email).
        Scan(&u.ID, &u.Name, &u.Email, &u.Password)
    
    if err == sql.ErrNoRows {
        log.Printf("User not found with email: %s", email)
        return nil, ErrUserNotFound
    }
    
    if err != nil {
        log.Printf("Error fetching user by email: %v", err)
    }
    
    return &u, err
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
    log.Printf("Fetching user by ID: %d", id)
    
    var u models.User
    err := r.db.QueryRow(ctx,
        "SELECT id, name, email FROM users WHERE id=$1", id).
        Scan(&u.ID, &u.Name, &u.Email)
    
    if err == sql.ErrNoRows {
        log.Printf("User not found with ID: %d", id)
        return nil, ErrUserNotFound
    }
    
    if err != nil {
        log.Printf("Error fetching user by ID: %v", err)
    }
    
    return &u, err
}

func (r *userRepository) GetAll(ctx context.Context) ([]models.User, error) {
    log.Printf("Fetching all users")
    
    rows, err := r.db.Query(ctx, "SELECT id, name, email FROM users ORDER BY id")
    if err != nil {
        log.Printf("Error fetching all users: %v", err)
        return nil, err
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var u models.User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
            log.Printf("Error scanning user row: %v", err)
            return nil, err
        }
        users = append(users, u)
    }

    log.Printf("Fetched %d users", len(users))
    return users, nil
}

func (r *userRepository) Update(ctx context.Context, u *models.User) error {
    log.Printf("Updating user ID: %d", u.ID)
    
    result, err := r.db.Exec(ctx, 
        "UPDATE users SET name=$1, email=$2 WHERE id=$3",
        u.Name, u.Email, u.ID)
    if err != nil {
        log.Printf("Error updating user: %v", err)
        return err
    }
    
    rowsAffected := result.RowsAffected()
    if rowsAffected == 0 {
        log.Printf("User not found for update: %d", u.ID)
        return ErrUserNotFound
    }
    
    log.Printf("User updated successfully: %d", u.ID)
    return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
    log.Printf("Deleting user ID: %d", id)
    
    result, err := r.db.Exec(ctx, "DELETE FROM users WHERE id=$1", id)
    if err != nil {
        log.Printf("Error deleting user: %v", err)
        return err
    }
    
    rowsAffected := result.RowsAffected()
    if rowsAffected == 0 {
        log.Printf("User not found for deletion: %d", id)
        return ErrUserNotFound
    }
    
    log.Printf("User deleted successfully: %d", id)
    return nil
}