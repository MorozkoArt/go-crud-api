package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "github.com/MorozkoArt/go-crud-api/internal/repository"
    "github.com/MorozkoArt/go-crud-api/internal/models"
    "github.com/MorozkoArt/go-crud-api/internal/auth"
    "github.com/MorozkoArt/go-crud-api/internal/middleware"
)

type Handler struct {
    repo      *repository.Repository
    jwtAuth   *auth.JWTService
}

func NewHandler(repo *repository.Repository, jwtAuth *auth.JWTService) *Handler {
    return &Handler{
        repo:    repo,
        jwtAuth: jwtAuth,
    }
}

type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

type AuthRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (h *Handler) RegisterRouter(r chi.Router) {
    r.Post("/register", h.Register)
    r.Post("/login", h.Login)
    
    r.Group(func(r chi.Router) {
        r.Use(middleware.AuthMiddleware(h.jwtAuth))
        
        r.Get("/", h.GetAllUser)
        r.Get("/{id}", h.GetUserByID)
        r.Put("/{id}", h.UpdateUser)
        r.Delete("/{id}", h.DeleteUser)
    })
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    var u models.User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        sendError(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.repo.Create(r.Context(), &u); err != nil {
        if err == repository.ErrUserExists {
            sendError(w, "User with this email already exists", http.StatusConflict)
        } else {
            sendError(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    sendSuccess(w, "User registered successfully", http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    var req AuthRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    user, err := h.repo.VerifyPassword(r.Context(), req.Email, req.Password)
    if err != nil {
        sendError(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    token, err := h.jwtAuth.GenerateToken(user.Id, user.Email)
    if err != nil {
        sendError(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "token": token,
        "user": map[string]interface{}{
            "id":    user.Id,
            "name":  user.Name,
            "email": user.Email,
        },
    }

    sendSuccess(w, response, http.StatusOK)
}

func (h *Handler) GetAllUser(w http.ResponseWriter, r *http.Request) {
    users, err := h.repo.GetAll(r.Context())
    if err != nil {
        sendError(w, err.Error(), http.StatusInternalServerError)
        return
    }
    sendSuccess(w, users, http.StatusOK)
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        sendError(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    u, err := h.repo.GetById(r.Context(), id)
    if err != nil {
        if err == repository.ErrUserNotFound {
            sendError(w, "User not found", http.StatusNotFound)
        } else {
            sendError(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    sendSuccess(w, u, http.StatusOK)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        sendError(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var u models.User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        sendError(w, err.Error(), http.StatusBadRequest)
        return
    }

    u.Id = id
    if err := h.repo.Update(r.Context(), &u); err != nil {
        if err == repository.ErrUserNotFound {
            sendError(w, "User not found", http.StatusNotFound)
        } else {
            sendError(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    sendSuccess(w, "User updated successfully", http.StatusOK)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        sendError(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    if err := h.repo.Delete(r.Context(), id); err != nil {
        if err == repository.ErrUserNotFound {
            sendError(w, "User not found", http.StatusNotFound)
        } else {
            sendError(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(Response{
        Success: false,
        Error:   message,
    })
}

func sendSuccess(w http.ResponseWriter, data interface{}, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(Response{
        Success: true,
        Data:    data,
    })
}

