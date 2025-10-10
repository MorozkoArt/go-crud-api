package main

import (
    "context"
    "fmt"
    "log"
    "net/http"

    "github.com/MorozkoArt/go-crud-api/internal/config"
    "github.com/MorozkoArt/go-crud-api/internal/db"
    "github.com/MorozkoArt/go-crud-api/internal/handlers"
    "github.com/MorozkoArt/go-crud-api/internal/repository"
    "github.com/MorozkoArt/go-crud-api/internal/services"
    "github.com/MorozkoArt/go-crud-api/internal/router"
)

func main() {
    ctx := context.Background()

    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Configuration loading error: %v", err)
    }

    pool, err := db.NewPostgresDB(ctx, cfg)
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }
    defer pool.Close()

    userRepo := repository.NewUserRepository(pool)
    authService := services.NewAuthService(cfg.Auth.JWTSecret, cfg.Auth.TokenExpiry)
    userService := services.NewUserService(userRepo, authService)
    userHandler := handlers.NewUserHandler(userService)

    r := router.NewRouter(userHandler, authService)

    addr := fmt.Sprintf(":%d", cfg.Server.Port)
    log.Printf("Server starting on %s", addr)
    
    if err := http.ListenAndServe(addr, r); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}