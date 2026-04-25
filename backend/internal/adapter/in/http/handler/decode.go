package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/kittitad/jeeb/internal/adapter/in/http/middleware"
	"github.com/kittitad/jeeb/internal/domain"
	"github.com/kittitad/jeeb/pkg/apperror"
)

var validate = validator.New()

func decodeAndValidate(r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return apperror.ValidationError(err.Error())
	}
	if err := validate.Struct(dst); err != nil {
		return apperror.ValidationError(err.Error())
	}
	return nil
}

func userIDFromCtx(r *http.Request) string {
	user, _ := r.Context().Value(middleware.UserCtxKey).(*domain.User)
	if user == nil {
		return ""
	}
	return user.ID
}
