package in

import (
	"context"

	"github.com/kittitad/jeeb/internal/domain"
	"github.com/kittitad/jeeb/pkg/pagination"
)

type CreateStudyRequest struct {
	Subject  string `validate:"required"`
	Duration int    `validate:"required,gt=0"`
	Notes    string
}

type UpdateStudyRequest struct {
	Subject  string `validate:"omitempty"`
	Duration int    `validate:"omitempty,gt=0"`
	Notes    string
}

type StudyStats struct {
	ThisWeek  int            `json:"this_week"`
	ThisMonth int            `json:"this_month"`
	Total     int64          `json:"total"`
	BySubject map[string]int `json:"by_subject"`
}

type StudyUseCase interface {
	Create(ctx context.Context, userID string, req CreateStudyRequest) (*domain.StudySession, error)
	GetByID(ctx context.Context, userID, id string) (*domain.StudySession, error)
	List(ctx context.Context, userID string, opts pagination.Params) ([]*domain.StudySession, *pagination.Meta, error)
	Update(ctx context.Context, userID, id string, req UpdateStudyRequest) (*domain.StudySession, error)
	Delete(ctx context.Context, userID, id string) error
	GetStats(ctx context.Context, userID string) (*StudyStats, error)
}
