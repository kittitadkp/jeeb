package in

import (
	"context"

	"github.com/kittitad/jeeb/internal/domain"
	"github.com/kittitad/jeeb/pkg/pagination"
)

type CreateWorkoutRequest struct {
	Type      domain.WorkoutType `validate:"required,oneof=strength cardio flexibility"`
	Duration  int                `validate:"required,gt=0"`
	Exercises []domain.Exercise  `validate:"omitempty,dive"`
	Notes     string
}

type UpdateWorkoutRequest struct {
	Type      domain.WorkoutType `validate:"omitempty,oneof=strength cardio flexibility"`
	Duration  int                `validate:"omitempty,gt=0"`
	Exercises []domain.Exercise
	Notes     string
}

type WorkoutStats struct {
	ThisWeek  int            `json:"this_week"`
	ThisMonth int            `json:"this_month"`
	Total     int64          `json:"total"`
	Streak    int            `json:"streak"`
	ByType    map[string]int `json:"by_type"`
}

type WorkoutUseCase interface {
	Create(ctx context.Context, userID string, req CreateWorkoutRequest) (*domain.Workout, error)
	GetByID(ctx context.Context, userID, id string) (*domain.Workout, error)
	List(ctx context.Context, userID string, opts pagination.Params) ([]*domain.Workout, *pagination.Meta, error)
	Update(ctx context.Context, userID, id string, req UpdateWorkoutRequest) (*domain.Workout, error)
	Delete(ctx context.Context, userID, id string) error
	GetStats(ctx context.Context, userID string) (*WorkoutStats, error)
}
