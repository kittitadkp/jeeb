package in

import (
	"context"

	"github.com/kittitadkp/jeeb-learning/internal/domain"
)

type TopicUseCase interface {
	List(ctx context.Context) ([]*domain.Topic, error)
	GetByID(ctx context.Context, id string) (*domain.Topic, error)
	Create(ctx context.Context, req CreateTopicRequest) (*domain.Topic, error)
	Update(ctx context.Context, id string, req UpdateTopicRequest) (*domain.Topic, error)
	Delete(ctx context.Context, id string) error
}

type CreateTopicRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Icon        string `json:"icon"`
}

type UpdateTopicRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Icon        string `json:"icon"`
}
