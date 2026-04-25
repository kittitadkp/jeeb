package out

import (
	"context"

	"github.com/kittitad/jeeb/internal/domain"
	"github.com/kittitad/jeeb/pkg/pagination"
)

type UserRepository interface {
	FindByKeycloakID(ctx context.Context, keycloakID string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Upsert(ctx context.Context, user *domain.User) error
}

type WorkoutRepository interface {
	Create(ctx context.Context, workout *domain.Workout) error
	FindByID(ctx context.Context, id string) (*domain.Workout, error)
	FindByUserID(ctx context.Context, userID string, opts pagination.Params) ([]*domain.Workout, int64, error)
	Update(ctx context.Context, workout *domain.Workout) error
	Delete(ctx context.Context, id string) error
}

type StudyRepository interface {
	Create(ctx context.Context, session *domain.StudySession) error
	FindByID(ctx context.Context, id string) (*domain.StudySession, error)
	FindByUserID(ctx context.Context, userID string, opts pagination.Params) ([]*domain.StudySession, int64, error)
	Update(ctx context.Context, session *domain.StudySession) error
	Delete(ctx context.Context, id string) error
}

type SleepRepository interface {
	Create(ctx context.Context, record *domain.SleepRecord) error
	FindByID(ctx context.Context, id string) (*domain.SleepRecord, error)
	FindByUserID(ctx context.Context, userID string, opts pagination.Params) ([]*domain.SleepRecord, int64, error)
	Update(ctx context.Context, record *domain.SleepRecord) error
	Delete(ctx context.Context, id string) error
}

type FinanceRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	FindByID(ctx context.Context, id string) (*domain.Transaction, error)
	FindByUserID(ctx context.Context, userID string, opts pagination.Params) ([]*domain.Transaction, int64, error)
	Update(ctx context.Context, tx *domain.Transaction) error
	Delete(ctx context.Context, id string) error
}

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	FindByID(ctx context.Context, id string) (*domain.Event, error)
	FindByUserID(ctx context.Context, userID string, opts pagination.Params) ([]*domain.Event, int64, error)
	Update(ctx context.Context, event *domain.Event) error
	Delete(ctx context.Context, id string) error
}
