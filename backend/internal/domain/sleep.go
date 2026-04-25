package domain

import "time"

type SleepRecord struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	StartTime time.Time `bson:"start_time" json:"start_time"`
	EndTime   time.Time `bson:"end_time" json:"end_time"`
	Quality   int       `bson:"quality" json:"quality"`
	Notes     string    `bson:"notes" json:"notes"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (s *SleepRecord) DurationMinutes() int {
	return int(s.EndTime.Sub(s.StartTime).Minutes())
}
