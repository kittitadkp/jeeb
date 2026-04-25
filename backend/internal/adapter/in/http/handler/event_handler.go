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

type EventHandler struct {
	uc in.EventUseCase
}

func NewEventHandler(uc in.EventUseCase) *EventHandler {
	return &EventHandler{uc: uc}
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string           `json:"title" validate:"required"`
		Type  domain.EventType `json:"type" validate:"required,oneof=workout study sleep finance custom"`
		Start time.Time        `json:"start" validate:"required"`
		End   time.Time        `json:"end" validate:"required"`
	}
	if err := decodeAndValidate(r, &req); err != nil {
		middleware.RespondError(w, err)
		return
	}

	event, err := h.uc.Create(r.Context(), userIDFromCtx(r), in.CreateEventRequest{
		Title: req.Title,
		Type:  req.Type,
		Start: req.Start,
		End:   req.End,
	})
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusCreated, event)
}

func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	opts := pagination.FromRequest(r)
	events, meta, err := h.uc.List(r.Context(), userIDFromCtx(r), opts)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"data": events, "meta": meta})
}

func (h *EventHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	event, err := h.uc.GetByID(r.Context(), userIDFromCtx(r), id)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, event)
}

func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Title string           `json:"title"`
		Type  domain.EventType `json:"type" validate:"omitempty,oneof=workout study sleep finance custom"`
		Start *time.Time       `json:"start"`
		End   *time.Time       `json:"end"`
	}
	if err := decodeAndValidate(r, &req); err != nil {
		middleware.RespondError(w, err)
		return
	}

	event, err := h.uc.Update(r.Context(), userIDFromCtx(r), id, in.UpdateEventRequest{
		Title: req.Title,
		Type:  req.Type,
		Start: req.Start,
		End:   req.End,
	})
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, event)
}

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.uc.Delete(r.Context(), userIDFromCtx(r), id); err != nil {
		middleware.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) SyncToCalendar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.uc.SyncToCalendar(r.Context(), userIDFromCtx(r), id); err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, map[string]string{"status": "synced"})
}
