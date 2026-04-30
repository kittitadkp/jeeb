package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/kittitadkp/jeeb-learning/internal/adapter/in/http/handler"
	"github.com/kittitadkp/jeeb-learning/internal/adapter/in/http/middleware"
)

type Handlers struct {
	User     *handler.UserHandler
	Topic    *handler.TopicHandler
	Item     *handler.ItemHandler
	Progress *handler.ProgressHandler
}

func NewRouter(h Handlers, auth *middleware.AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		middleware.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.Authenticate)

		r.Get("/me", h.User.Me)

		r.Get("/stats", h.Progress.GetStats)

		r.Route("/topics", func(r chi.Router) {
			r.Get("/", h.Topic.List)
			r.Post("/", h.Topic.Create)
			r.Get("/{id}", h.Topic.GetByID)
			r.Put("/{id}", h.Topic.Update)
			r.Delete("/{id}", h.Topic.Delete)

			r.Get("/{id}/items", h.Item.List)
			r.Post("/{id}/items", h.Item.Create)
			r.Put("/{id}/items/{itemId}", h.Item.Update)
			r.Delete("/{id}/items/{itemId}", h.Item.Delete)

			r.Get("/{id}/progress", h.Progress.GetTopicProgress)
			r.Delete("/{id}/progress", h.Progress.ResetTopic)
		})

		r.Put("/progress/{itemId}", h.Progress.Upsert)
	})

	return r
}
