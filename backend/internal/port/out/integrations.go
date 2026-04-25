package out

import (
	"context"

	"github.com/kittitad/jeeb/internal/domain"
)

type CalendarPort interface {
	CreateEvent(ctx context.Context, userID string, event *domain.Event) (externalID string, err error)
	UpdateEvent(ctx context.Context, userID string, event *domain.Event) error
	DeleteEvent(ctx context.Context, userID, externalID string) error
}

type NotificationPort interface {
	Send(ctx context.Context, userID string, message string) error
}
