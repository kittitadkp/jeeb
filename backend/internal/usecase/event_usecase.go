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

type eventUseCase struct {
	repo     out.EventRepository
	calendar out.CalendarPort
}

func NewEventUseCase(repo out.EventRepository, calendar out.CalendarPort) in.EventUseCase {
	return &eventUseCase{repo: repo, calendar: calendar}
}

func (uc *eventUseCase) Create(ctx context.Context, userID string, req in.CreateEventRequest) (*domain.Event, error) {
	now := time.Now()
	event := &domain.Event{
		UserID:    userID,
		Title:     req.Title,
		Type:      req.Type,
		Start:     req.Start,
		End:       req.End,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.repo.Create(ctx, event); err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to create event", http.StatusInternalServerError)
	}
	return event, nil
}

func (uc *eventUseCase) GetByID(ctx context.Context, userID, id string) (*domain.Event, error) {
	event, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	if event.UserID != userID {
		return nil, apperror.ErrForbidden
	}
	return event, nil
}

func (uc *eventUseCase) List(ctx context.Context, userID string, opts pagination.Params) ([]*domain.Event, *pagination.Meta, error) {
	slog.Debug("listing events", "userID", userID, "page", opts.Page, "limit", opts.Limit)

	events, total, err := uc.repo.FindByUserID(ctx, userID, opts)
	if err != nil {
		return nil, nil, apperror.New(apperror.CodeInternalError, "failed to list events", http.StatusInternalServerError)
	}
	return events, pagination.NewMeta(opts.Page, opts.Limit, total), nil
}

func (uc *eventUseCase) Update(ctx context.Context, userID, id string, req in.UpdateEventRequest) (*domain.Event, error) {
	event, err := uc.GetByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		event.Title = req.Title
	}
	if req.Type != "" {
		event.Type = req.Type
	}
	if req.Start != nil {
		event.Start = *req.Start
	}
	if req.End != nil {
		event.End = *req.End
	}
	event.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, event); err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to update event", http.StatusInternalServerError)
	}
	return event, nil
}

func (uc *eventUseCase) Delete(ctx context.Context, userID, id string) error {
	event, err := uc.GetByID(ctx, userID, id)
	if err != nil {
		return err
	}
	if event.ExternalID != "" && uc.calendar != nil {
		_ = uc.calendar.DeleteEvent(ctx, userID, event.ExternalID)
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		return apperror.New(apperror.CodeInternalError, "failed to delete event", http.StatusInternalServerError)
	}
	return nil
}

func (uc *eventUseCase) SyncToCalendar(ctx context.Context, userID, id string) error {
	if uc.calendar == nil {
		return apperror.New(apperror.CodeInternalError, "calendar integration not configured", 503)
	}

	event, err := uc.GetByID(ctx, userID, id)
	if err != nil {
		return err
	}

	if event.ExternalID != "" {
		if err := uc.calendar.UpdateEvent(ctx, userID, event); err != nil {
			return apperror.New(apperror.CodeInternalError, "failed to update calendar event", http.StatusInternalServerError)
		}
		return nil
	}

	externalID, err := uc.calendar.CreateEvent(ctx, userID, event)
	if err != nil {
		return apperror.New(apperror.CodeInternalError, "failed to create calendar event", http.StatusInternalServerError)
	}

	event.ExternalID = externalID
	event.UpdatedAt = time.Now()
	if err := uc.repo.Update(ctx, event); err != nil {
		return apperror.New(apperror.CodeInternalError, "failed to update event after calendar sync", http.StatusInternalServerError)
	}
	return nil
}
