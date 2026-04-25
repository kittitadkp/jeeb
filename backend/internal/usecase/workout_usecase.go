package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/kittitad/jeeb/internal/domain"
	"github.com/kittitad/jeeb/internal/port/in"
	"github.com/kittitad/jeeb/internal/port/out"
	"github.com/kittitad/jeeb/pkg/apperror"
	"github.com/kittitad/jeeb/pkg/pagination"
)

type workoutUseCase struct {
	repo out.WorkoutRepository
}

func NewWorkoutUseCase(repo out.WorkoutRepository) in.WorkoutUseCase {
	return &workoutUseCase{repo: repo}
}

func (uc *workoutUseCase) Create(ctx context.Context, userID string, req in.CreateWorkoutRequest) (*domain.Workout, error) {
	now := time.Now()
	workout := &domain.Workout{
		UserID:    userID,
		Type:      req.Type,
		Duration:  req.Duration,
		Exercises: req.Exercises,
		Notes:     req.Notes,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.repo.Create(ctx, workout); err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to create workout", http.StatusInternalServerError)
	}
	return workout, nil
}

func (uc *workoutUseCase) GetByID(ctx context.Context, userID, id string) (*domain.Workout, error) {
	workout, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	if workout.UserID != userID {
		return nil, apperror.ErrForbidden
	}
	return workout, nil
}

func (uc *workoutUseCase) List(ctx context.Context, userID string, opts pagination.Params) ([]*domain.Workout, *pagination.Meta, error) {
	workouts, total, err := uc.repo.FindByUserID(ctx, userID, opts)
	if err != nil {
		return nil, nil, apperror.New(apperror.CodeInternalError, "failed to list workouts", http.StatusInternalServerError)
	}
	return workouts, pagination.NewMeta(opts.Page, opts.Limit, total), nil
}

func (uc *workoutUseCase) Update(ctx context.Context, userID, id string, req in.UpdateWorkoutRequest) (*domain.Workout, error) {
	workout, err := uc.GetByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	if req.Type != "" {
		workout.Type = req.Type
	}
	if req.Duration > 0 {
		workout.Duration = req.Duration
	}
	if req.Exercises != nil {
		workout.Exercises = req.Exercises
	}
	workout.Notes = req.Notes
	workout.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, workout); err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to update workout", http.StatusInternalServerError)
	}
	return workout, nil
}

func (uc *workoutUseCase) Delete(ctx context.Context, userID, id string) error {
	if _, err := uc.GetByID(ctx, userID, id); err != nil {
		return err
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		return apperror.New(apperror.CodeInternalError, "failed to delete workout", http.StatusInternalServerError)
	}
	return nil
}

func (uc *workoutUseCase) GetStats(ctx context.Context, userID string) (*in.WorkoutStats, error) {
	// Fetch all to compute stats — in production, push aggregation to MongoDB
	opts := pagination.Params{Page: 1, Limit: 1000, Sort: "-created_at"}
	workouts, total, err := uc.repo.FindByUserID(ctx, userID, opts)
	if err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to get workout stats", http.StatusInternalServerError)
	}

	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	stats := &in.WorkoutStats{
		Total:  total,
		ByType: make(map[string]int),
	}

	dateSet := make(map[string]bool)
	for _, w := range workouts {
		stats.ByType[string(w.Type)]++
		if w.CreatedAt.After(weekStart) {
			stats.ThisWeek++
		}
		if w.CreatedAt.After(monthStart) {
			stats.ThisMonth++
		}
		dateSet[w.CreatedAt.Format("2006-01-02")] = true
	}

	// Compute streak: count consecutive days backwards from today (or yesterday if no workout today)
	startOffset := 0
	if !dateSet[now.Format("2006-01-02")] {
		startOffset = 1
	}
	for i := startOffset; ; i++ {
		d := now.AddDate(0, 0, -i).Format("2006-01-02")
		if dateSet[d] {
			stats.Streak++
		} else {
			break
		}
	}

	return stats, nil
}
