package services

import (
    "context"
    "errors"
    "log"

    "github.com/MorozkoArt/go-crud-api/internal/models"
    "github.com/MorozkoArt/go-crud-api/internal/repository"
    "github.com/MorozkoArt/go-crud-api/internal/utils"
)

type UserService interface {
    Register(ctx context.Context, req *models.RegisterRequest) error
    Login(ctx context.Context, req *models.LoginRequest) (*models.UserResponse, string, error)
    GetAllUsers(ctx context.Context) ([]models.UserResponse, error)
    GetUserByID(ctx context.Context, id int64) (*models.UserResponse, error)
    UpdateUser(ctx context.Context, id int64, req *models.UpdateUserRequest) error
    DeleteUser(ctx context.Context, id int64) error
}

type userService struct {
    userRepo    repository.UserRepository
    authService AuthService
}

func NewUserService(userRepo repository.UserRepository, authService AuthService) UserService {
    return &userService{
        userRepo:    userRepo,
        authService: authService,
    }
}

func (s *userService) Register(ctx context.Context, req *models.RegisterRequest) error {
    log.Printf("Service: Registering user: %s", req.Email)
    
    user := &models.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
    }
    
    return s.userRepo.Create(ctx, user)
}

func (s *userService) Login(ctx context.Context, req *models.LoginRequest) (*models.UserResponse, string, error) {
    log.Printf("Service: Login attempt for: %s", req.Email)
    
    user, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err != nil {
        log.Printf("Service: Login failed - user not found: %s", req.Email)
        return nil, "", errors.New("invalid credentials")
    }

    if !utils.CheckPasswordHash(req.Password, user.Password) {
        log.Printf("Service: Login failed - invalid password for: %s", req.Email)
        return nil, "", errors.New("invalid credentials")
    }

    token, err := s.authService.GenerateToken(user.ID, user.Email)
    if err != nil {
        log.Printf("Service: Token generation failed: %v", err)
        return nil, "", err
    }

    log.Printf("Service: Login successful for: %s", req.Email)
    return &models.UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }, token, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]models.UserResponse, error) {
    log.Printf("Service: Fetching all users")
    
    users, err := s.userRepo.GetAll(ctx)
    if err != nil {
        return nil, err
    }

    var response []models.UserResponse
    for _, user := range users {
        response = append(response, models.UserResponse{
            ID:    user.ID,
            Name:  user.Name,
            Email: user.Email,
        })
    }

    return response, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*models.UserResponse, error) {
    log.Printf("Service: Fetching user by ID: %d", id)
    
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    return &models.UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int64, req *models.UpdateUserRequest) error {
    log.Printf("Service: Updating user ID: %d", id)
    
    user := &models.User{
        ID:    id,
        Name:  req.Name,
        Email: req.Email,
    }
    
    return s.userRepo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
    log.Printf("Service: Deleting user ID: %d", id)
    return s.userRepo.Delete(ctx, id)
}