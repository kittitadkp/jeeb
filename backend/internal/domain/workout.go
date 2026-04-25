package domain

import "time"

type WorkoutType string

const (
	WorkoutStrength    WorkoutType = "strength"
	WorkoutCardio      WorkoutType = "cardio"
	WorkoutFlexibility WorkoutType = "flexibility"
)

type Exercise struct {
	Name   string  `bson:"name" json:"name"`
	Sets   int     `bson:"sets" json:"sets"`
	Reps   int     `bson:"reps" json:"reps"`
	Weight float64 `bson:"weight" json:"weight"`
}

type Workout struct {
	ID        string      `bson:"_id,omitempty" json:"id"`
	UserID    string      `bson:"user_id" json:"user_id"`
	Type      WorkoutType `bson:"type" json:"type"`
	Duration  int         `bson:"duration" json:"duration"`
	Exercises []Exercise  `bson:"exercises" json:"exercises"`
	Notes     string      `bson:"notes" json:"notes"`
	CreatedAt time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time   `bson:"updated_at" json:"updated_at"`
}
