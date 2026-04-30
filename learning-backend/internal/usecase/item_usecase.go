package usecase

import (
	"context"
	"time"

	"github.com/kittitadkp/jeeb-learning/internal/domain"
	"github.com/kittitadkp/jeeb-learning/internal/port/in"
	"github.com/kittitadkp/jeeb-learning/internal/port/out"
	"github.com/kittitadkp/jeeb-learning/pkg/apperror"
	"github.com/kittitadkp/jeeb-learning/pkg/pagination"
)

type itemUseCase struct {
	repo out.ItemRepository
}

func NewItemUseCase(repo out.ItemRepository) in.ItemUseCase {
	return &itemUseCase{repo: repo}
}

func (uc *itemUseCase) List(ctx context.Context, topicID string, opts pagination.Params, category string) ([]*domain.Item, *pagination.Meta, error) {
	items, total, err := uc.repo.FindByTopicID(ctx, topicID, opts, category)
	if err != nil {
		return nil, nil, apperror.ErrInternal
	}
	return items, pagination.NewMeta(opts.Page, opts.Limit, total), nil
}

func (uc *itemUseCase) GetByID(ctx context.Context, id string) (*domain.Item, error) {
	item, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	return item, nil
}

func (uc *itemUseCase) Create(ctx context.Context, topicID string, req in.CreateItemRequest) (*domain.Item, error) {
	now := time.Now()
	item := &domain.Item{
		TopicID:   topicID,
		Term:      req.Term,
		Meaning:   req.Meaning,
		Example:   req.Example,
		Hint:      req.Hint,
		Category:  req.Category,
		SortOrder: req.SortOrder,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.repo.Create(ctx, item); err != nil {
		return nil, apperror.ErrInternal
	}
	return item, nil
}

func (uc *itemUseCase) Update(ctx context.Context, id string, req in.UpdateItemRequest) (*domain.Item, error) {
	item, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	if req.Term != "" {
		item.Term = req.Term
	}
	if req.Meaning != "" {
		item.Meaning = req.Meaning
	}
	if req.Example != "" {
		item.Example = req.Example
	}
	if req.Hint != "" {
		item.Hint = req.Hint
	}
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.SortOrder != 0 {
		item.SortOrder = req.SortOrder
	}
	item.UpdatedAt = time.Now()
	if err := uc.repo.Update(ctx, item); err != nil {
		return nil, apperror.ErrInternal
	}
	return item, nil
}

func (uc *itemUseCase) Delete(ctx context.Context, id string) error {
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		return apperror.ErrNotFound
	}
	return uc.repo.Delete(ctx, id)
}
