package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/kittitad/jeeb/internal/port/in"
	"github.com/kittitad/jeeb/pkg/apperror"
)

type contextKey string

const UserCtxKey contextKey = "user"

type Claims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type AuthMiddleware struct {
	verifier    *oidc.IDTokenVerifier
	userUseCase in.UserUseCase
}

func NewAuthMiddleware(verifier *oidc.IDTokenVerifier, userUseCase in.UserUseCase) *AuthMiddleware {
	return &AuthMiddleware{verifier: verifier, userUseCase: userUseCase}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r)
		if token == "" {
			RespondError(w, apperror.ErrUnauthorized)
			return
		}

		idToken, err := m.verifier.Verify(r.Context(), token)
		if err != nil {
			RespondError(w, apperror.ErrUnauthorized)
			return
		}

		var claims Claims
		if err := idToken.Claims(&claims); err != nil {
			slog.Debug("Authenticating Claims", "error", err)
			RespondError(w, apperror.ErrUnauthorized)
			return
		}

		user, err := m.userUseCase.GetOrCreate(r.Context(), claims.Sub, claims.Email, claims.Name)
		if err != nil {
			slog.Debug("Authenticating GetOrCreate", "error", err)
			RespondError(w, apperror.ErrInternal)
			return
		}

		ctx := context.WithValue(r.Context(), UserCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(auth, "Bearer ")
}
