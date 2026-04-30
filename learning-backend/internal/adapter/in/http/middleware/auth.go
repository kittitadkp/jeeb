package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/kittitadkp/jeeb-learning/internal/domain"
	"github.com/kittitadkp/jeeb-learning/internal/port/in"
	"github.com/kittitadkp/jeeb-learning/pkg/apperror"
)

type contextKey string

const UserCtxKey contextKey = "user"

type Claims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type AuthMiddleware struct {
	verifier     *oidc.IDTokenVerifier
	userUseCase  in.UserUseCase
	upstreamAuth bool
}

func NewAuthMiddleware(verifier *oidc.IDTokenVerifier, userUseCase in.UserUseCase, upstreamAuth bool) *AuthMiddleware {
	return &AuthMiddleware{verifier: verifier, userUseCase: userUseCase, upstreamAuth: upstreamAuth}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r)
		if token == "" {
			RespondError(w, apperror.ErrUnauthorized)
			return
		}

		var claims *Claims
		if m.upstreamAuth {
			// Kong already verified the JWT — just decode the payload.
			var err error
			claims, err = decodeJWTPayload(token)
			if err != nil {
				slog.Debug("jwt decode error", "error", err)
				RespondError(w, apperror.ErrUnauthorized)
				return
			}
		} else {
			idToken, err := m.verifier.Verify(r.Context(), token)
			if err != nil {
				RespondError(w, apperror.ErrUnauthorized)
				return
			}
			claims = &Claims{}
			if err := idToken.Claims(claims); err != nil {
				slog.Debug("auth claims error", "error", err)
				RespondError(w, apperror.ErrUnauthorized)
				return
			}
		}

		user, err := m.userUseCase.GetOrCreate(r.Context(), claims.Sub, claims.Email, claims.Name)
		if err != nil {
			RespondError(w, apperror.ErrInternal)
			return
		}

		ctx := context.WithValue(r.Context(), UserCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// decodeJWTPayload extracts claims from the JWT payload without signature verification.
// Only safe behind a gateway (Kong) that has already verified the token.
func decodeJWTPayload(tokenStr string) (*Claims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}
	if claims.Sub == "" {
		return nil, errors.New("missing sub claim")
	}
	return &claims, nil
}

func extractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(auth, "Bearer ")
}

func UserFromCtx(r *http.Request) *domain.User {
	user, _ := r.Context().Value(UserCtxKey).(*domain.User)
	return user
}
