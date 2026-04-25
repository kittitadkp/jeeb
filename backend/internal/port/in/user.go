package in

import (
	"context"

	"github.com/kittitad/jeeb/internal/domain"
)

type UserUseCase interface {
	GetOrCreate(ctx context.Context, keycloakID, email, displayName string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}
