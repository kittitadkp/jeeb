# Workout Module

## Domain

```go
type Workout struct {
    ID        string
    UserID    string
    Type      WorkoutType  // strength, cardio, flexibility
    Duration  time.Duration
    Exercises []Exercise
    CreatedAt time.Time
}
```

## Use Cases

- CreateWorkout
- ListWorkouts
- GetWorkoutStats
- SyncToCalendar

## Repository Interface

```go
type WorkoutRepository interface {
    Save(ctx context.Context, w *Workout) error
    FindByUserID(ctx context.Context, userID string) ([]*Workout, error)
    FindByID(ctx context.Context, id string) (*Workout, error)
}
```
