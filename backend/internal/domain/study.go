package domain

import "time"

type StudySession struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	Subject   string    `bson:"subject" json:"subject"`
	Duration  int       `bson:"duration" json:"duration"`
	Notes     string    `bson:"notes" json:"notes"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
