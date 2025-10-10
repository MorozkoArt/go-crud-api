package services

import (
    "time"

    "github.com/MorozkoArt/go-crud-api/internal/utils"
)

type AuthService interface {
    GenerateToken(userID int64, email string) (string, error)
    ValidateToken(tokenString string) (*utils.Claims, error)
}

type authService struct {
    jwtService *utils.JWTService
}

func NewAuthService(secretKey string, expiry time.Duration) AuthService {
    return &authService{
        jwtService: utils.NewJWTService(secretKey, expiry),
    }
}

func (s *authService) GenerateToken(userID int64, email string) (string, error) {
    return s.jwtService.GenerateToken(userID, email)
}

func (s *authService) ValidateToken(tokenString string) (*utils.Claims, error) {
    return s.jwtService.ValidateToken(tokenString)
}