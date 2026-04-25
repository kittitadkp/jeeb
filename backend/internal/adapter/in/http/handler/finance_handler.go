package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/kittitad/jeeb/internal/adapter/in/http/middleware"
	"github.com/kittitad/jeeb/internal/domain"
	"github.com/kittitad/jeeb/internal/port/in"
	"github.com/kittitad/jeeb/pkg/pagination"
)

type FinanceHandler struct {
	uc in.FinanceUseCase
}

func NewFinanceHandler(uc in.FinanceUseCase) *FinanceHandler {
	return &FinanceHandler{uc: uc}
}

func (h *FinanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type     domain.TransactionType `json:"type" validate:"required,oneof=income expense"`
		Amount   float64                `json:"amount" validate:"required,gt=0"`
		Category string                 `json:"category" validate:"required"`
		Date     time.Time              `json:"date" validate:"required"`
		Notes    string                 `json:"notes"`
	}
	if err := decodeAndValidate(r, &req); err != nil {
		middleware.RespondError(w, err)
		return
	}

	tx, err := h.uc.Create(r.Context(), userIDFromCtx(r), in.CreateTransactionRequest{
		Type:     req.Type,
		Amount:   req.Amount,
		Category: req.Category,
		Date:     req.Date,
		Notes:    req.Notes,
	})
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusCreated, tx)
}

func (h *FinanceHandler) List(w http.ResponseWriter, r *http.Request) {
	opts := pagination.FromRequest(r)
	txs, meta, err := h.uc.List(r.Context(), userIDFromCtx(r), opts)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"data": txs, "meta": meta})
}

func (h *FinanceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tx, err := h.uc.GetByID(r.Context(), userIDFromCtx(r), id)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, tx)
}

func (h *FinanceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Type     domain.TransactionType `json:"type" validate:"omitempty,oneof=income expense"`
		Amount   float64                `json:"amount" validate:"omitempty,gt=0"`
		Category string                 `json:"category"`
		Date     *time.Time             `json:"date"`
		Notes    string                 `json:"notes"`
	}
	if err := decodeAndValidate(r, &req); err != nil {
		middleware.RespondError(w, err)
		return
	}

	tx, err := h.uc.Update(r.Context(), userIDFromCtx(r), id, in.UpdateTransactionRequest{
		Type:     req.Type,
		Amount:   req.Amount,
		Category: req.Category,
		Date:     req.Date,
		Notes:    req.Notes,
	})
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, tx)
}

func (h *FinanceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.uc.Delete(r.Context(), userIDFromCtx(r), id); err != nil {
		middleware.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *FinanceHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetStats(r.Context(), userIDFromCtx(r))
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, stats)
}

func (h *FinanceHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.uc.GetCategories(r.Context(), userIDFromCtx(r))
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"data": categories})
}
