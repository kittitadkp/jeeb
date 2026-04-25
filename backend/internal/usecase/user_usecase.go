package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/kittitad/jeeb/internal/domain"
	"github.com/kittitad/jeeb/internal/port/out"
	"github.com/kittitad/jeeb/pkg/apperror"
)

type userUseCase struct {
	repo out.UserRepository
}

func NewUserUseCase(repo out.UserRepository) *userUseCase {
	return &userUseCase{repo: repo}
}

func (uc *userUseCase) GetOrCreate(ctx context.Context, keycloakID, email, displayName string) (*domain.User, error) {
	user, err := uc.repo.FindByKeycloakID(ctx, keycloakID)
	if err == nil {
		return user, nil
	}

	now := time.Now()
	user = &domain.User{
		KeycloakID:  keycloakID,
		Email:       email,
		DisplayName: displayName,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := uc.repo.Upsert(ctx, user); err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to create or update user", http.StatusInternalServerError)
	}
	return user, nil
}

func (uc *userUseCase) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to find user", http.StatusInternalServerError)
	}
	return user, nil
}
