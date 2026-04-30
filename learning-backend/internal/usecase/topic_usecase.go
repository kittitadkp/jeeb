package usecase

import (
	"context"
	"time"

	"github.com/kittitadkp/jeeb-learning/internal/domain"
	"github.com/kittitadkp/jeeb-learning/internal/port/in"
	"github.com/kittitadkp/jeeb-learning/internal/port/out"
	"github.com/kittitadkp/jeeb-learning/pkg/apperror"
)

type topicUseCase struct {
	repo out.TopicRepository
}

func NewTopicUseCase(repo out.TopicRepository) in.TopicUseCase {
	return &topicUseCase{repo: repo}
}

func (uc *topicUseCase) List(ctx context.Context) ([]*domain.Topic, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *topicUseCase) GetByID(ctx context.Context, id string) (*domain.Topic, error) {
	topic, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	return topic, nil
}

func (uc *topicUseCase) Create(ctx context.Context, req in.CreateTopicRequest) (*domain.Topic, error) {
	now := time.Now()
	topic := &domain.Topic{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Icon:        req.Icon,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := uc.repo.Create(ctx, topic); err != nil {
		return nil, apperror.ErrInternal
	}
	return topic, nil
}

func (uc *topicUseCase) Update(ctx context.Context, id string, req in.UpdateTopicRequest) (*domain.Topic, error) {
	topic, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	if req.Name != "" {
		topic.Name = req.Name
	}
	if req.Description != "" {
		topic.Description = req.Description
	}
	if req.Category != "" {
		topic.Category = req.Category
	}
	if req.Icon != "" {
		topic.Icon = req.Icon
	}
	topic.UpdatedAt = time.Now()
	if err := uc.repo.Update(ctx, topic); err != nil {
		return nil, apperror.ErrInternal
	}
	return topic, nil
}

func (uc *topicUseCase) Delete(ctx context.Context, id string) error {
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		return apperror.ErrNotFound
	}
	return uc.repo.Delete(ctx, id)
}
