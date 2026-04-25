package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kittitad/jeeb/internal/adapter/in/http/middleware"
	"github.com/kittitad/jeeb/internal/domain"
	"github.com/kittitad/jeeb/internal/port/in"
	"github.com/kittitad/jeeb/pkg/pagination"
)

type WorkoutHandler struct {
	uc in.WorkoutUseCase
}

func NewWorkoutHandler(uc in.WorkoutUseCase) *WorkoutHandler {
	return &WorkoutHandler{uc: uc}
}

func (h *WorkoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type      domain.WorkoutType `json:"type" validate:"required,oneof=strength cardio flexibility"`
		Duration  int                `json:"duration" validate:"required,gt=0"`
		Exercises []domain.Exercise  `json:"exercises"`
		Notes     string             `json:"notes"`
	}
	if err := decodeAndValidate(r, &req); err != nil {
		middleware.RespondError(w, err)
		return
	}

	workout, err := h.uc.Create(r.Context(), userIDFromCtx(r), in.CreateWorkoutRequest{
		Type:      req.Type,
		Duration:  req.Duration,
		Exercises: req.Exercises,
		Notes:     req.Notes,
	})
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusCreated, workout)
}

func (h *WorkoutHandler) List(w http.ResponseWriter, r *http.Request) {
	opts := pagination.FromRequest(r)
	workouts, meta, err := h.uc.List(r.Context(), userIDFromCtx(r), opts)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"data": workouts, "meta": meta})
}

func (h *WorkoutHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	workout, err := h.uc.GetByID(r.Context(), userIDFromCtx(r), id)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, workout)
}

func (h *WorkoutHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Type      domain.WorkoutType `json:"type" validate:"omitempty,oneof=strength cardio flexibility"`
		Duration  int                `json:"duration" validate:"omitempty,gt=0"`
		Exercises []domain.Exercise  `json:"exercises"`
		Notes     string             `json:"notes"`
	}
	if err := decodeAndValidate(r, &req); err != nil {
		middleware.RespondError(w, err)
		return
	}

	workout, err := h.uc.Update(r.Context(), userIDFromCtx(r), id, in.UpdateWorkoutRequest{
		Type:      req.Type,
		Duration:  req.Duration,
		Exercises: req.Exercises,
		Notes:     req.Notes,
	})
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, workout)
}

func (h *WorkoutHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.uc.Delete(r.Context(), userIDFromCtx(r), id); err != nil {
		middleware.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkoutHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetStats(r.Context(), userIDFromCtx(r))
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, stats)
}
