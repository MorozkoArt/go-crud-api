package router

import (
    "github.com/go-chi/chi/v5"
    "github.com/MorozkoArt/go-crud-api/internal/handlers"
    "github.com/MorozkoArt/go-crud-api/internal/middleware"
    "github.com/MorozkoArt/go-crud-api/internal/services"
)

func NewRouter(userHandler *handlers.UserHandler, authService services.AuthService) *chi.Mux {
    r := chi.NewRouter()
    
    r.Use(middleware.Logger)
    
    r.Route("/api/users", func(r chi.Router) {
        r.Post("/register", userHandler.Register)
        r.Post("/login", userHandler.Login)
        
        r.Group(func(r chi.Router) {
            r.Use(middleware.AuthMiddleware(authService))
            
            r.Get("/", userHandler.GetAllUsers)
            r.Get("/{id}", userHandler.GetUserByID)
            r.Put("/{id}", userHandler.UpdateUser)
            r.Delete("/{id}", userHandler.DeleteUser)
        })
    })
    
    return r
}