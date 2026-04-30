package apperror

import "net/http"

type Code string

const (
	CodeUnauthorized  Code = "UNAUTHORIZED"
	CodeForbidden     Code = "FORBIDDEN"
	CodeNotFound      Code = "NOT_FOUND"
	CodeValidation    Code = "VALIDATION_ERROR"
	CodeConflict      Code = "CONFLICT"
	CodeInternalError Code = "INTERNAL_ERROR"
)

type AppError struct {
	Code    Code        `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	status  int
}

func (e *AppError) Error() string   { return e.Message }
func (e *AppError) HTTPStatus() int { return e.status }

func New(code Code, message string, status int) *AppError {
	return &AppError{Code: code, Message: message, status: status}
}

func NewWithDetails(code Code, message string, status int, details interface{}) *AppError {
	return &AppError{Code: code, Message: message, status: status, Details: details}
}

var (
	ErrUnauthorized = New(CodeUnauthorized, "unauthorized", http.StatusUnauthorized)
	ErrForbidden    = New(CodeForbidden, "forbidden", http.StatusForbidden)
	ErrNotFound     = New(CodeNotFound, "not found", http.StatusNotFound)
	ErrInternal     = New(CodeInternalError, "internal server error", http.StatusInternalServerError)
	ErrConflict     = New(CodeConflict, "resource already exists", http.StatusConflict)
)

func ValidationError(details interface{}) *AppError {
	return NewWithDetails(CodeValidation, "invalid request", http.StatusBadRequest, details)
}
