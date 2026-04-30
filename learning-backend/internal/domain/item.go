package domain

import "time"

type Item struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	TopicID   string    `bson:"topic_id" json:"topic_id"`
	Term      string    `bson:"term" json:"term"`
	Meaning   string    `bson:"meaning" json:"meaning"`
	Example   string    `bson:"example" json:"example"`
	Hint      string    `bson:"hint" json:"hint"`
	Category  string    `bson:"category" json:"category"`
	SortOrder int       `bson:"sort_order" json:"sort_order"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
