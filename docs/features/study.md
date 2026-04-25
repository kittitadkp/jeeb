# Study Module

## Domain

```go
type StudySession struct {
    ID        string
    UserID    string
    Subject   string
    Duration  time.Duration
    Notes     string
    CreatedAt time.Time
}
```

## Use Cases

- CreateStudySession
- ListStudySessions
- GetStudyStats
- SetStudyReminder

## Tracking

- Total hours per subject
- Daily/weekly goals
- Streak tracking
