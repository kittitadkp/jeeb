package middleware

import (
	"log/slog"
	"net/http"

	"github.com/kittitad/jeeb/pkg/apperror"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered", "error", rec)
				RespondError(w, apperror.ErrInternal)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
