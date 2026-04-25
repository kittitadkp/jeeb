package usecase

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/kittitad/jeeb/internal/domain"
	"github.com/kittitad/jeeb/internal/port/in"
	"github.com/kittitad/jeeb/internal/port/out"
	"github.com/kittitad/jeeb/pkg/apperror"
	"github.com/kittitad/jeeb/pkg/pagination"
)

type studyUseCase struct {
	repo out.StudyRepository
}

func NewStudyUseCase(repo out.StudyRepository) in.StudyUseCase {
	return &studyUseCase{repo: repo}
}

func (uc *studyUseCase) Create(ctx context.Context, userID string, req in.CreateStudyRequest) (*domain.StudySession, error) {
	now := time.Now()
	session := &domain.StudySession{
		UserID:    userID,
		Subject:   req.Subject,
		Duration:  req.Duration,
		Notes:     req.Notes,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.repo.Create(ctx, session); err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to create study session", http.StatusInternalServerError)
	}
	return session, nil
}

func (uc *studyUseCase) GetByID(ctx context.Context, userID, id string) (*domain.StudySession, error) {
	session, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	if session.UserID != userID {
		return nil, apperror.ErrForbidden
	}
	return session, nil
}

func (uc *studyUseCase) List(ctx context.Context, userID string, opts pagination.Params) ([]*domain.StudySession, *pagination.Meta, error) {
	slog.Debug("listing study sessions", "userID", userID, "page", opts.Page, "limit", opts.Limit)

	sessions, total, err := uc.repo.FindByUserID(ctx, userID, opts)
	if err != nil {
		return nil, nil, apperror.New(apperror.CodeInternalError, "failed to list study sessions", http.StatusInternalServerError)
	}
	return sessions, pagination.NewMeta(opts.Page, opts.Limit, total), nil
}

func (uc *studyUseCase) Update(ctx context.Context, userID, id string, req in.UpdateStudyRequest) (*domain.StudySession, error) {
	session, err := uc.GetByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	if req.Subject != "" {
		session.Subject = req.Subject
	}
	if req.Duration > 0 {
		session.Duration = req.Duration
	}
	session.Notes = req.Notes
	session.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, session); err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to update study session", http.StatusInternalServerError)
	}
	return session, nil
}

func (uc *studyUseCase) Delete(ctx context.Context, userID, id string) error {
	if _, err := uc.GetByID(ctx, userID, id); err != nil {
		return err
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		return apperror.New(apperror.CodeInternalError, "failed to delete study session", http.StatusInternalServerError)
	}
	return nil
}

func (uc *studyUseCase) GetStats(ctx context.Context, userID string) (*in.StudyStats, error) {
	opts := pagination.Params{Page: 1, Limit: 1000, Sort: "-created_at"}
	sessions, total, err := uc.repo.FindByUserID(ctx, userID, opts)
	if err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to get study stats", http.StatusInternalServerError)
	}

	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	stats := &in.StudyStats{
		Total:     total,
		BySubject: make(map[string]int),
	}

	for _, s := range sessions {
		stats.BySubject[s.Subject]++
		if s.CreatedAt.After(weekStart) {
			stats.ThisWeek += s.Duration
		}
		if s.CreatedAt.After(monthStart) {
			stats.ThisMonth += s.Duration
		}
	}

	return stats, nil
}
