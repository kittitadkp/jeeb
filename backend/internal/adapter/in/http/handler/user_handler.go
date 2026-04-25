package handler

import (
	"net/http"

	"github.com/kittitad/jeeb/internal/adapter/in/http/middleware"
	"github.com/kittitad/jeeb/internal/domain"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler { return &UserHandler{} }

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(middleware.UserCtxKey).(*domain.User)
	middleware.RespondJSON(w, http.StatusOK, user)
}
