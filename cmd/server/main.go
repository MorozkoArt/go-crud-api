package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MorozkoArt/go-crud-api/internal/config"
	"github.com/MorozkoArt/go-crud-api/internal/db"
	"github.com/MorozkoArt/go-crud-api/internal/user"
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

	// testUser := &user.User{
	// 	Name: 	  "Alex",
	// 	Email: 	  "alex@gmail.com",
	// 	Password: "123456",
	// }

	// if err := repo.Create(ctx, testUser); err != nil {
	// 	log.Println("Error creating user")
	// } else {
	// 	fmt.Println("User added!")
	// }

	users, _ := repo.GetAll(ctx)
	fmt.Println("Users in the database:", users)

	repo.Delete(ctx, 1)

	users, _ = repo.GetAll(ctx)
	fmt.Println("Users in the database:", users)

}
