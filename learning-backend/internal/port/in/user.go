package in

import (
	"context"

	"github.com/kittitadkp/jeeb-learning/internal/domain"
)

type UserUseCase interface {
	GetOrCreate(ctx context.Context, keycloakID, email, displayName string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}
