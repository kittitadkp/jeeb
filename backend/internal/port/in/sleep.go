package in

import (
	"context"
	"time"

	"github.com/kittitad/jeeb/internal/domain"
	"github.com/kittitad/jeeb/pkg/pagination"
)

type CreateSleepRequest struct {
	StartTime time.Time `validate:"required"`
	EndTime   time.Time `validate:"required,gtfield=StartTime"`
	Quality   int       `validate:"required,min=1,max=5"`
	Notes     string
}

type UpdateSleepRequest struct {
	StartTime *time.Time `validate:"omitempty"`
	EndTime   *time.Time `validate:"omitempty"`
	Quality   int        `validate:"omitempty,min=1,max=5"`
	Notes     string
}

type SleepStats struct {
	ThisWeek       int     `json:"this_week"`
	ThisMonth      int     `json:"this_month"`
	AvgDuration    float64 `json:"avg_duration_minutes"`
	AvgQuality     float64 `json:"avg_quality"`
}

type SleepUseCase interface {
	Create(ctx context.Context, userID string, req CreateSleepRequest) (*domain.SleepRecord, error)
	GetByID(ctx context.Context, userID, id string) (*domain.SleepRecord, error)
	List(ctx context.Context, userID string, opts pagination.Params) ([]*domain.SleepRecord, *pagination.Meta, error)
	Update(ctx context.Context, userID, id string, req UpdateSleepRequest) (*domain.SleepRecord, error)
	Delete(ctx context.Context, userID, id string) error
	GetStats(ctx context.Context, userID string) (*SleepStats, error)
}
