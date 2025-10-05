package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/MorozkoArt/go-crud-api/internal/config"
	"github.com/MorozkoArt/go-crud-api/internal/db"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	pool, err := db.NewPostgresDB(ctx, cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer pool.Close()

	port := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Сервер запущен на порту %s", port)

	http.ListenAndServe(port, nil)
}
