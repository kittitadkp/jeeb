package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kittitadkp/jeeb-learning/pkg/apperror"
)

func RespondJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func RespondError(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		RespondJSON(w, appErr.HTTPStatus(), map[string]interface{}{"error": appErr})
		return
	}
	RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
		"error": apperror.ErrInternal,
	})
}
