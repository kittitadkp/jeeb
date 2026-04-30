package in

import (
	"context"

	"github.com/kittitadkp/jeeb-learning/internal/domain"
	"github.com/kittitadkp/jeeb-learning/pkg/pagination"
)

type ItemUseCase interface {
	List(ctx context.Context, topicID string, opts pagination.Params, category string) ([]*domain.Item, *pagination.Meta, error)
	GetByID(ctx context.Context, id string) (*domain.Item, error)
	Create(ctx context.Context, topicID string, req CreateItemRequest) (*domain.Item, error)
	Update(ctx context.Context, id string, req UpdateItemRequest) (*domain.Item, error)
	Delete(ctx context.Context, id string) error
}

type CreateItemRequest struct {
	Term      string `json:"term" validate:"required"`
	Meaning   string `json:"meaning" validate:"required"`
	Example   string `json:"example"`
	Hint      string `json:"hint"`
	Category  string `json:"category"`
	SortOrder int    `json:"sort_order"`
}

type UpdateItemRequest struct {
	Term      string `json:"term"`
	Meaning   string `json:"meaning"`
	Example   string `json:"example"`
	Hint      string `json:"hint"`
	Category  string `json:"category"`
	SortOrder int    `json:"sort_order"`
}
