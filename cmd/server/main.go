package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MorozkoArt/go-crud-api/internal/config"
	"github.com/MorozkoArt/go-crud-api/internal/db"
	"github.com/MorozkoArt/go-crud-api/internal/user"
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

	repo := user.NewRepository(pool)
	handler := router.NewHandler(repo)

	r := chi.NewRouter()
	r.Route("/users", handler.RegisterRouter)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("Сервер запущен на %s\n", addr)
	http.ListenAndServe(addr, r)

}
