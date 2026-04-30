package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kittitadkp/jeeb-learning/internal/adapter/in/http/middleware"
	"github.com/kittitadkp/jeeb-learning/internal/port/in"
)

type TopicHandler struct {
	uc in.TopicUseCase
}

func NewTopicHandler(uc in.TopicUseCase) *TopicHandler {
	return &TopicHandler{uc: uc}
}

func (h *TopicHandler) List(w http.ResponseWriter, r *http.Request) {
	topics, err := h.uc.List(r.Context())
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{"data": topics})
}

func (h *TopicHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	topic, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, topic)
}

func (h *TopicHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req in.CreateTopicRequest
	if err := decodeAndValidate(r, &req); err != nil {
		middleware.RespondError(w, err)
		return
	}
	topic, err := h.uc.Create(r.Context(), req)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusCreated, topic)
}

func (h *TopicHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req in.UpdateTopicRequest
	if err := decodeAndValidate(r, &req); err != nil {
		middleware.RespondError(w, err)
		return
	}
	topic, err := h.uc.Update(r.Context(), id, req)
	if err != nil {
		middleware.RespondError(w, err)
		return
	}
	middleware.RespondJSON(w, http.StatusOK, topic)
}

func (h *TopicHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.uc.Delete(r.Context(), id); err != nil {
		middleware.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
