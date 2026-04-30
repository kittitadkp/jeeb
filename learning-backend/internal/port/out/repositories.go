package out

import (
	"context"

	"github.com/kittitadkp/jeeb-learning/internal/domain"
	"github.com/kittitadkp/jeeb-learning/pkg/pagination"
)

type UserRepository interface {
	FindByKeycloakID(ctx context.Context, keycloakID string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Upsert(ctx context.Context, user *domain.User) error
}

type TopicRepository interface {
	Create(ctx context.Context, topic *domain.Topic) error
	FindAll(ctx context.Context) ([]*domain.Topic, error)
	FindByID(ctx context.Context, id string) (*domain.Topic, error)
	Update(ctx context.Context, topic *domain.Topic) error
	Delete(ctx context.Context, id string) error
}

type ItemRepository interface {
	Create(ctx context.Context, item *domain.Item) error
	InsertMany(ctx context.Context, items []*domain.Item) error
	FindByTopicID(ctx context.Context, topicID string, opts pagination.Params, category string) ([]*domain.Item, int64, error)
	FindAllByTopicID(ctx context.Context, topicID string) ([]*domain.Item, error)
	FindByID(ctx context.Context, id string) (*domain.Item, error)
	Update(ctx context.Context, item *domain.Item) error
	Delete(ctx context.Context, id string) error
}

type ProgressRepository interface {
	Upsert(ctx context.Context, p *domain.UserProgress) error
	FindByUserAndTopic(ctx context.Context, userID, topicID string) ([]*domain.UserProgress, error)
	FindByUserAndItem(ctx context.Context, userID, itemID string) (*domain.UserProgress, error)
	DeleteByUserAndTopic(ctx context.Context, userID, topicID string) error
}
