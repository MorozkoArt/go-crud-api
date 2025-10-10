package main

import (
    "context"
    "fmt"
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"

    "github.com/MorozkoArt/go-crud-api/internal/config"
    "github.com/MorozkoArt/go-crud-api/internal/db"
    "github.com/MorozkoArt/go-crud-api/internal/repository"
    "github.com/MorozkoArt/go-crud-api/internal/handlers"
    "github.com/MorozkoArt/go-crud-api/internal/auth"
    "github.com/MorozkoArt/go-crud-api/internal/middleware"
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

    jwtService := auth.NewJWTService(cfg.Auth.JWTSecret, cfg.Auth.TokenExpiry)

    repo := repository.NewRepository(pool)
    handler := handlers.NewHandler(repo, jwtService)

    r := chi.NewRouter()
    
    r.Use(middleware.Logger)
    
    r.Route("/api/users", handler.RegisterRouter)

    addr := fmt.Sprintf(":%d", cfg.Server.Port)
    fmt.Printf("Server started on %s\n", addr)
    
    log.Fatal(http.ListenAndServe(addr, r))
}