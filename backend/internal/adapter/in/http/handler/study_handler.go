package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kittitad/jeeb/internal/adapter/in/http/middleware"
	"github.com/kittitad/jeeb/internal/port/in"
	"github.com/kittitad/jeeb/pkg/pagination"
)

type StudyHandler struct {
	uc in.StudyUseCase
}

func NewStudyHandler(uc in.StudyUseCase) *StudyHandler {
	return &StudyHandler{uc: uc}
}

func (h *StudyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Subject  string `json:"subject" validate:"required"`
		Duration int    `json:"duration" validate:"required,gt=0"`
		Notes    string `json:"notes"`
	}
	if err := decodeAndValidate(r, &req); err != nil {
		middleware.RespondError(w, err)
		return
	}

	session, err := h.uc.Create(r.Context(), userIDFromCtx(r), in.CreateStudyRequest{
		Subject:  req.Subject,
		Duration: req.Duration,
		Notes:    req.Notes,
	})
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusCreated, session)
}

func (h *StudyHandler) List(w http.ResponseWriter, r *http.Request) {
	opts := pagination.FromRequest(r)
	sessions, meta, err := h.uc.List(r.Context(), userIDFromCtx(r), opts)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"data": sessions, "meta": meta})
}

func (h *StudyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	session, err := h.uc.GetByID(r.Context(), userIDFromCtx(r), id)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, session)
}

func (h *StudyHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Subject  string `json:"subject"`
		Duration int    `json:"duration" validate:"omitempty,gt=0"`
		Notes    string `json:"notes"`
	}
	if err := decodeAndValidate(r, &req); err != nil {
		middleware.RespondError(w, err)
		return
	}

	session, err := h.uc.Update(r.Context(), userIDFromCtx(r), id, in.UpdateStudyRequest{
		Subject:  req.Subject,
		Duration: req.Duration,
		Notes:    req.Notes,
	})
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, session)
}

func (h *StudyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.uc.Delete(r.Context(), userIDFromCtx(r), id); err != nil {
		middleware.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *StudyHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetStats(r.Context(), userIDFromCtx(r))
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, stats)
}
