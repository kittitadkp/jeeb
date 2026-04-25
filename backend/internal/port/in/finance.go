package in

import (
	"context"
	"time"

	"github.com/kittitad/jeeb/internal/domain"
	"github.com/kittitad/jeeb/pkg/pagination"
)

type CreateTransactionRequest struct {
	Type     domain.TransactionType `validate:"required,oneof=income expense"`
	Amount   float64                `validate:"required,gt=0"`
	Category string                 `validate:"required"`
	Date     time.Time              `validate:"required"`
	Notes    string
}

type UpdateTransactionRequest struct {
	Type     domain.TransactionType `validate:"omitempty,oneof=income expense"`
	Amount   float64                `validate:"omitempty,gt=0"`
	Category string                 `validate:"omitempty"`
	Date     *time.Time
	Notes    string
}

type FinanceStats struct {
	ThisMonth struct {
		Income  float64 `json:"income"`
		Expense float64 `json:"expense"`
		Net     float64 `json:"net"`
	} `json:"this_month"`
	ByCategory map[string]float64 `json:"by_category"`
}

type FinanceUseCase interface {
	Create(ctx context.Context, userID string, req CreateTransactionRequest) (*domain.Transaction, error)
	GetByID(ctx context.Context, userID, id string) (*domain.Transaction, error)
	List(ctx context.Context, userID string, opts pagination.Params) ([]*domain.Transaction, *pagination.Meta, error)
	Update(ctx context.Context, userID, id string, req UpdateTransactionRequest) (*domain.Transaction, error)
	Delete(ctx context.Context, userID, id string) error
	GetStats(ctx context.Context, userID string) (*FinanceStats, error)
	GetCategories(ctx context.Context, userID string) ([]string, error)
}
