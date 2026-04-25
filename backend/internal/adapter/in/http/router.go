package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/kittitad/jeeb/internal/adapter/in/http/handler"
	"github.com/kittitad/jeeb/internal/adapter/in/http/middleware"
)

type Handlers struct {
	User    *handler.UserHandler
	Workout *handler.WorkoutHandler
	Study   *handler.StudyHandler
	Sleep   *handler.SleepHandler
	Finance *handler.FinanceHandler
	Event   *handler.EventHandler
}

func NewRouter(h Handlers, auth *middleware.AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
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

		r.Route("/workouts", func(r chi.Router) {
			r.Get("/", h.Workout.List)
			r.Post("/", h.Workout.Create)
			r.Get("/stats", h.Workout.GetStats)
			r.Get("/{id}", h.Workout.GetByID)
			r.Put("/{id}", h.Workout.Update)
			r.Delete("/{id}", h.Workout.Delete)
		})

		r.Route("/study", func(r chi.Router) {
			r.Get("/", h.Study.List)
			r.Post("/", h.Study.Create)
			r.Get("/stats", h.Study.GetStats)
			r.Get("/{id}", h.Study.GetByID)
			r.Put("/{id}", h.Study.Update)
			r.Delete("/{id}", h.Study.Delete)
		})

		r.Route("/sleep", func(r chi.Router) {
			r.Get("/", h.Sleep.List)
			r.Post("/", h.Sleep.Create)
			r.Get("/stats", h.Sleep.GetStats)
			r.Get("/{id}", h.Sleep.GetByID)
			r.Put("/{id}", h.Sleep.Update)
			r.Delete("/{id}", h.Sleep.Delete)
		})

		r.Route("/finance", func(r chi.Router) {
			r.Get("/", h.Finance.List)
			r.Post("/", h.Finance.Create)
			r.Get("/stats", h.Finance.GetStats)
			r.Get("/categories", h.Finance.GetCategories)
			r.Get("/{id}", h.Finance.GetByID)
			r.Put("/{id}", h.Finance.Update)
			r.Delete("/{id}", h.Finance.Delete)
		})

		r.Route("/events", func(r chi.Router) {
			r.Get("/", h.Event.List)
			r.Post("/", h.Event.Create)
			r.Get("/{id}", h.Event.GetByID)
			r.Put("/{id}", h.Event.Update)
			r.Delete("/{id}", h.Event.Delete)
			r.Post("/{id}/sync", h.Event.SyncToCalendar)
		})
	})

	return r
}
