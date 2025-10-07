package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/MorozkoArt/go-crud-api/internal/user"
)

type Handler struct {
	repo *user.Repository
}

func NewHandler(repo *user.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRouter(r chi.Router) {
	r.Post("/", h.CreateUser)
	r.Get("/", h.GetAllUser)
	r.Get("/{id}", h.GetUserByID)
	r.Put("/{id}", h.UpdateUser)
	r.Delete("/{id}", h.DeleteUser)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.repo.Create(r.Context(), &u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}


