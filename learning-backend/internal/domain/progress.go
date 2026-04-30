package domain

import "time"

const (
	StatusLearning = "learning"
	StatusMastered = "mastered"
)

type UserProgress struct {
	ID             string    `bson:"_id,omitempty" json:"id"`
	UserID         string    `bson:"user_id" json:"user_id"`
	TopicID        string    `bson:"topic_id" json:"topic_id"`
	ItemID         string    `bson:"item_id" json:"item_id"`
	Status         string    `bson:"status" json:"status"`
	ReviewCount    int       `bson:"review_count" json:"review_count"`
	LastReviewedAt time.Time `bson:"last_reviewed_at" json:"last_reviewed_at"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at" json:"updated_at"`
}
