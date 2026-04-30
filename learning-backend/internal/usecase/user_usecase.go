package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/kittitadkp/jeeb-learning/internal/domain"
	"github.com/kittitadkp/jeeb-learning/internal/port/out"
	"github.com/kittitadkp/jeeb-learning/pkg/apperror"
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
		return nil, apperror.New(apperror.CodeInternalError, "failed to create user", http.StatusInternalServerError)
	}
	return user, nil
}

func (uc *userUseCase) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	return user, nil
}
