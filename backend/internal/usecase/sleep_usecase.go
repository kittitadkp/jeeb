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

type sleepUseCase struct {
	repo out.SleepRepository
}

func NewSleepUseCase(repo out.SleepRepository) in.SleepUseCase {
	return &sleepUseCase{repo: repo}
}

func (uc *sleepUseCase) Create(ctx context.Context, userID string, req in.CreateSleepRequest) (*domain.SleepRecord, error) {
	now := time.Now()
	record := &domain.SleepRecord{
		UserID:    userID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Quality:   req.Quality,
		Notes:     req.Notes,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.repo.Create(ctx, record); err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to create sleep record", http.StatusInternalServerError)
	}
	return record, nil
}

func (uc *sleepUseCase) GetByID(ctx context.Context, userID, id string) (*domain.SleepRecord, error) {
	record, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	if record.UserID != userID {
		return nil, apperror.ErrForbidden
	}
	return record, nil
}

func (uc *sleepUseCase) List(ctx context.Context, userID string, opts pagination.Params) ([]*domain.SleepRecord, *pagination.Meta, error) {
	records, total, err := uc.repo.FindByUserID(ctx, userID, opts)
	if err != nil {
		return nil, nil, apperror.New(apperror.CodeInternalError, "failed to list sleep records", http.StatusInternalServerError)
	}
	return records, pagination.NewMeta(opts.Page, opts.Limit, total), nil
}

func (uc *sleepUseCase) Update(ctx context.Context, userID, id string, req in.UpdateSleepRequest) (*domain.SleepRecord, error) {
	record, err := uc.GetByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	if req.StartTime != nil {
		record.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		record.EndTime = *req.EndTime
	}
	if req.Quality > 0 {
		record.Quality = req.Quality
	}
	record.Notes = req.Notes
	record.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, record); err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to update sleep record", http.StatusInternalServerError)
	}
	return record, nil
}

func (uc *sleepUseCase) Delete(ctx context.Context, userID, id string) error {
	if _, err := uc.GetByID(ctx, userID, id); err != nil {
		return err
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		return apperror.New(apperror.CodeInternalError, "failed to delete sleep record", http.StatusInternalServerError)
	}
	return nil
}

func (uc *sleepUseCase) GetStats(ctx context.Context, userID string) (*in.SleepStats, error) {
	opts := pagination.Params{Page: 1, Limit: 1000, Sort: "-created_at"}
	records, _, err := uc.repo.FindByUserID(ctx, userID, opts)
	if err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to get sleep stats", http.StatusInternalServerError)
	}

	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	stats := &in.SleepStats{}
	var totalDuration, totalQuality float64
	count := 0

	for _, r := range records {
		if r.CreatedAt.After(weekStart) {
			stats.ThisWeek++
		}
		if r.CreatedAt.After(monthStart) {
			stats.ThisMonth++
		}
		totalDuration += float64(r.DurationMinutes())
		totalQuality += float64(r.Quality)
		count++
	}

	if count > 0 {
		stats.AvgDuration = totalDuration / float64(count)
		stats.AvgQuality = totalQuality / float64(count)
	}

	return stats, nil
}
