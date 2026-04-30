package handler

import (
	"net/http"

	"github.com/kittitadkp/jeeb-learning/internal/adapter/in/http/middleware"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler { return &UserHandler{} }

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	user := userFromCtx(r)
	middleware.RespondJSON(w, http.StatusOK, user)
}
