package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kittitadkp/jeeb-learning/internal/adapter/in/http/middleware"
	"github.com/kittitadkp/jeeb-learning/internal/port/in"
	"github.com/kittitadkp/jeeb-learning/pkg/pagination"
)

type ItemHandler struct {
	uc in.ItemUseCase
}

func NewItemHandler(uc in.ItemUseCase) *ItemHandler {
	return &ItemHandler{uc: uc}
}

func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	topicID := chi.URLParam(r, "id")
	opts := pagination.FromRequest(r)
	category := r.URL.Query().Get("category")

	items, meta, err := h.uc.List(r.Context(), topicID, opts, category)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"data": items, "meta": meta})
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	topicID := chi.URLParam(r, "id")
	var req in.CreateItemRequest
	if err := decodeAndValidate(r, &req); err != nil {
		middleware.RespondError(w, err)
		return
	}
	item, err := h.uc.Create(r.Context(), topicID, req)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusCreated, item)
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "itemId")
	var req in.UpdateItemRequest
	if err := decodeAndValidate(r, &req); err != nil {
		middleware.RespondError(w, err)
		return
	}
	item, err := h.uc.Update(r.Context(), itemID, req)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "itemId")
	if err := h.uc.Delete(r.Context(), itemID); err != nil {
		middleware.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
