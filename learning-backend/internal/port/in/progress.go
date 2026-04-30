package in

import (
	"context"

	"github.com/kittitadkp/jeeb-learning/internal/domain"
)

type ProgressUseCase interface {
	GetTopicProgress(ctx context.Context, userID, topicID string) (map[string]string, error)
	Upsert(ctx context.Context, userID, topicID, itemID, status string) (*domain.UserProgress, error)
	ResetTopic(ctx context.Context, userID, topicID string) error
	GetStats(ctx context.Context, userID string) ([]*TopicStats, error)
}

type TopicStats struct {
	TopicID  string `json:"topic_id"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Total    int    `json:"total"`
	Mastered int    `json:"mastered"`
	Learning int    `json:"learning"`
}
