package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kittitadkp/jeeb-learning/internal/adapter/in/http/middleware"
	"github.com/kittitadkp/jeeb-learning/internal/port/in"
)

type ProgressHandler struct {
	uc     in.ProgressUseCase
	itemUC in.ItemUseCase
}

func NewProgressHandler(uc in.ProgressUseCase, itemUC in.ItemUseCase) *ProgressHandler {
	return &ProgressHandler{uc: uc, itemUC: itemUC}
}

func (h *ProgressHandler) GetTopicProgress(w http.ResponseWriter, r *http.Request) {
	topicID := chi.URLParam(r, "id")
	userID := userIDFromCtx(r)

	progress, err := h.uc.GetTopicProgress(r.Context(), userID, topicID)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"data": progress})
}

func (h *ProgressHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "itemId")
	userID := userIDFromCtx(r)

	var req struct {
		TopicID string `json:"topic_id" validate:"required"`
		Status  string `json:"status" validate:"required"`
	}
	if err := decodeAndValidate(r, &req); err != nil {
		middleware.RespondError(w, err)
		return
	}

	p, err := h.uc.Upsert(r.Context(), userID, req.TopicID, itemID, req.Status)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, p)
}

func (h *ProgressHandler) ResetTopic(w http.ResponseWriter, r *http.Request) {
	topicID := chi.URLParam(r, "id")
	userID := userIDFromCtx(r)

	if err := h.uc.ResetTopic(r.Context(), userID, topicID); err != nil {
		middleware.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProgressHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromCtx(r)
	stats, err := h.uc.GetStats(r.Context(), userID)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"data": stats})
}
