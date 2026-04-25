package in

import (
	"context"
	"time"

	"github.com/kittitad/jeeb/internal/domain"
	"github.com/kittitad/jeeb/pkg/pagination"
)

type CreateEventRequest struct {
	Title string            `validate:"required"`
	Type  domain.EventType  `validate:"required,oneof=workout study sleep finance custom"`
	Start time.Time         `validate:"required"`
	End   time.Time         `validate:"required,gtfield=Start"`
}

type UpdateEventRequest struct {
	Title string           `validate:"omitempty"`
	Type  domain.EventType `validate:"omitempty,oneof=workout study sleep finance custom"`
	Start *time.Time
	End   *time.Time
}

type EventUseCase interface {
	Create(ctx context.Context, userID string, req CreateEventRequest) (*domain.Event, error)
	GetByID(ctx context.Context, userID, id string) (*domain.Event, error)
	List(ctx context.Context, userID string, opts pagination.Params) ([]*domain.Event, *pagination.Meta, error)
	Update(ctx context.Context, userID, id string, req UpdateEventRequest) (*domain.Event, error)
	Delete(ctx context.Context, userID, id string) error
	SyncToCalendar(ctx context.Context, userID, id string) error
}
