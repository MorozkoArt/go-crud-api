package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "github.com/MorozkoArt/go-crud-api/internal/models"
    "github.com/MorozkoArt/go-crud-api/internal/services"
    "github.com/MorozkoArt/go-crud-api/internal/utils"
)

type UserHandler struct {
    userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req models.RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if err := utils.ValidateStruct(req); err != nil {
        sendError(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.userService.Register(r.Context(), &req); err != nil {
        if err.Error() == "user already exists" {
            sendError(w, "User with this email already exists", http.StatusConflict)
        } else {
            sendError(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    sendSuccess(w, "User registered successfully", http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req models.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if err := utils.ValidateStruct(req); err != nil {
        sendError(w, err.Error(), http.StatusBadRequest)
        return
    }

    user, token, err := h.userService.Login(r.Context(), &req)
    if err != nil {
        sendError(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    response := map[string]interface{}{
        "token": token,
        "user":  user,
    }

    sendSuccess(w, response, http.StatusOK)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.userService.GetAllUsers(r.Context())
    if err != nil {
        sendError(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    sendSuccess(w, users, http.StatusOK)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        sendError(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    user, err := h.userService.GetUserByID(r.Context(), id)
    if err != nil {
        if err.Error() == "user not found" {
            sendError(w, "User not found", http.StatusNotFound)
        } else {
            sendError(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    sendSuccess(w, user, http.StatusOK)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        sendError(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var req models.UpdateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if err := utils.ValidateStruct(req); err != nil {
        sendError(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.userService.UpdateUser(r.Context(), id, &req); err != nil {
        if err.Error() == "user not found" {
            sendError(w, "User not found", http.StatusNotFound)
        } else {
            sendError(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    sendSuccess(w, "User updated successfully", http.StatusOK)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        sendError(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    if err := h.userService.DeleteUser(r.Context(), id); err != nil {
        if err.Error() == "user not found" {
            sendError(w, "User not found", http.StatusNotFound)
        } else {
            sendError(w, "Internal server error", http.StatusInternalServerError)
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