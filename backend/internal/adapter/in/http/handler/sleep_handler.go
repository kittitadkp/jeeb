package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/kittitad/jeeb/internal/adapter/in/http/middleware"
	"github.com/kittitad/jeeb/internal/port/in"
	"github.com/kittitad/jeeb/pkg/pagination"
)

type SleepHandler struct {
	uc in.SleepUseCase
}

func NewSleepHandler(uc in.SleepUseCase) *SleepHandler {
	return &SleepHandler{uc: uc}
}

func (h *SleepHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StartTime time.Time `json:"start_time" validate:"required"`
		EndTime   time.Time `json:"end_time" validate:"required"`
		Quality   int       `json:"quality" validate:"required,min=1,max=5"`
		Notes     string    `json:"notes"`
	}
	if err := decodeAndValidate(r, &req); err != nil {
		middleware.RespondError(w, err)
		return
	}

	record, err := h.uc.Create(r.Context(), userIDFromCtx(r), in.CreateSleepRequest{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Quality:   req.Quality,
		Notes:     req.Notes,
	})
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusCreated, record)
}

func (h *SleepHandler) List(w http.ResponseWriter, r *http.Request) {
	opts := pagination.FromRequest(r)
	records, meta, err := h.uc.List(r.Context(), userIDFromCtx(r), opts)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"data": records, "meta": meta})
}

func (h *SleepHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	record, err := h.uc.GetByID(r.Context(), userIDFromCtx(r), id)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, record)
}

func (h *SleepHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		StartTime *time.Time `json:"start_time"`
		EndTime   *time.Time `json:"end_time"`
		Quality   int        `json:"quality" validate:"omitempty,min=1,max=5"`
		Notes     string     `json:"notes"`
	}
	if err := decodeAndValidate(r, &req); err != nil {
		middleware.RespondError(w, err)
		return
	}

	record, err := h.uc.Update(r.Context(), userIDFromCtx(r), id, in.UpdateSleepRequest{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Quality:   req.Quality,
		Notes:     req.Notes,
	})
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, record)
}

func (h *SleepHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.uc.Delete(r.Context(), userIDFromCtx(r), id); err != nil {
		middleware.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *SleepHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetStats(r.Context(), userIDFromCtx(r))
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, stats)
}
