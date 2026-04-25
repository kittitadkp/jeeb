package domain

import "time"

type EventType string

const (
	EventWorkout EventType = "workout"
	EventStudy   EventType = "study"
	EventSleep   EventType = "sleep"
	EventFinance EventType = "finance"
	EventCustom  EventType = "custom"
)

type Event struct {
	ID         string    `bson:"_id,omitempty" json:"id"`
	UserID     string    `bson:"user_id" json:"user_id"`
	Title      string    `bson:"title" json:"title"`
	Type       EventType `bson:"type" json:"type"`
	Start      time.Time `bson:"start" json:"start"`
	End        time.Time `bson:"end" json:"end"`
	ExternalID string    `bson:"external_id" json:"external_id"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}
