# Sleep Module

## Domain

```go
type SleepRecord struct {
    ID        string
    UserID    string
    StartTime time.Time
    EndTime   time.Time
    Quality   int  // 1-5
    Notes     string
}
```

## Use Cases

- LogSleep
- GetSleepHistory
- GetSleepStats
- SetBedtimeReminder

## Metrics

- Average duration
- Sleep quality trends
- Consistency score
