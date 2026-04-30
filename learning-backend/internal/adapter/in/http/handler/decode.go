package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/kittitadkp/jeeb-learning/internal/adapter/in/http/middleware"
	"github.com/kittitadkp/jeeb-learning/internal/domain"
	"github.com/kittitadkp/jeeb-learning/pkg/apperror"
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

func userFromCtx(r *http.Request) *domain.User {
	return middleware.UserFromCtx(r)
}

func userIDFromCtx(r *http.Request) string {
	user := middleware.UserFromCtx(r)
	if user == nil {
		return ""
	}
	return user.ID
}
