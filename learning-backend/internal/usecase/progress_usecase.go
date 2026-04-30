package usecase

import (
	"context"
	"time"

	"github.com/kittitadkp/jeeb-learning/internal/domain"
	"github.com/kittitadkp/jeeb-learning/internal/port/in"
	"github.com/kittitadkp/jeeb-learning/internal/port/out"
	"github.com/kittitadkp/jeeb-learning/pkg/apperror"
)

type progressUseCase struct {
	progressRepo out.ProgressRepository
	topicRepo    out.TopicRepository
	itemRepo     out.ItemRepository
}

func NewProgressUseCase(progressRepo out.ProgressRepository, topicRepo out.TopicRepository, itemRepo out.ItemRepository) in.ProgressUseCase {
	return &progressUseCase{progressRepo: progressRepo, topicRepo: topicRepo, itemRepo: itemRepo}
}

func (uc *progressUseCase) GetTopicProgress(ctx context.Context, userID, topicID string) (map[string]string, error) {
	records, err := uc.progressRepo.FindByUserAndTopic(ctx, userID, topicID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	result := make(map[string]string, len(records))
	for _, r := range records {
		result[r.ItemID] = r.Status
	}
	return result, nil
}

func (uc *progressUseCase) Upsert(ctx context.Context, userID, topicID, itemID, status string) (*domain.UserProgress, error) {
	if status != domain.StatusLearning && status != domain.StatusMastered {
		return nil, apperror.ValidationError("status must be 'learning' or 'mastered'")
	}

	now := time.Now()
	existing, err := uc.progressRepo.FindByUserAndItem(ctx, userID, itemID)
	if err != nil {
		existing = &domain.UserProgress{
			UserID:    userID,
			TopicID:   topicID,
			ItemID:    itemID,
			CreatedAt: now,
		}
	}

	existing.Status = status
	existing.ReviewCount++
	existing.LastReviewedAt = now
	existing.UpdatedAt = now

	if err := uc.progressRepo.Upsert(ctx, existing); err != nil {
		return nil, apperror.ErrInternal
	}
	return existing, nil
}

func (uc *progressUseCase) ResetTopic(ctx context.Context, userID, topicID string) error {
	return uc.progressRepo.DeleteByUserAndTopic(ctx, userID, topicID)
}

func (uc *progressUseCase) GetStats(ctx context.Context, userID string) ([]*in.TopicStats, error) {
	topics, err := uc.topicRepo.FindAll(ctx)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	stats := make([]*in.TopicStats, 0, len(topics))
	for _, t := range topics {
		items, err := uc.itemRepo.FindAllByTopicID(ctx, t.ID)
		if err != nil {
			continue
		}
		progress, err := uc.progressRepo.FindByUserAndTopic(ctx, userID, t.ID)
		if err != nil {
			continue
		}

		mastered, learning := 0, 0
		for _, p := range progress {
			switch p.Status {
			case domain.StatusMastered:
				mastered++
			case domain.StatusLearning:
				learning++
			}
		}

		stats = append(stats, &in.TopicStats{
			TopicID:  t.ID,
			Name:     t.Name,
			Icon:     t.Icon,
			Total:    len(items),
			Mastered: mastered,
			Learning: learning,
		})
	}
	return stats, nil
}
