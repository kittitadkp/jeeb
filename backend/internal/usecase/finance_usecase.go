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

type financeUseCase struct {
	repo out.FinanceRepository
}

func NewFinanceUseCase(repo out.FinanceRepository) in.FinanceUseCase {
	return &financeUseCase{repo: repo}
}

func (uc *financeUseCase) Create(ctx context.Context, userID string, req in.CreateTransactionRequest) (*domain.Transaction, error) {
	now := time.Now()
	tx := &domain.Transaction{
		UserID:    userID,
		Type:      req.Type,
		Amount:    req.Amount,
		Category:  req.Category,
		Date:      req.Date,
		Notes:     req.Notes,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.repo.Create(ctx, tx); err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to create transaction", http.StatusInternalServerError)
	}
	return tx, nil
}

func (uc *financeUseCase) GetByID(ctx context.Context, userID, id string) (*domain.Transaction, error) {
	tx, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	if tx.UserID != userID {
		return nil, apperror.ErrForbidden
	}
	return tx, nil
}

func (uc *financeUseCase) List(ctx context.Context, userID string, opts pagination.Params) ([]*domain.Transaction, *pagination.Meta, error) {
	txs, total, err := uc.repo.FindByUserID(ctx, userID, opts)
	if err != nil {
		return nil, nil, apperror.New(apperror.CodeInternalError, "failed to list transactions", http.StatusInternalServerError)
	}
	return txs, pagination.NewMeta(opts.Page, opts.Limit, total), nil
}

func (uc *financeUseCase) Update(ctx context.Context, userID, id string, req in.UpdateTransactionRequest) (*domain.Transaction, error) {
	tx, err := uc.GetByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	if req.Type != "" {
		tx.Type = req.Type
	}
	if req.Amount > 0 {
		tx.Amount = req.Amount
	}
	if req.Category != "" {
		tx.Category = req.Category
	}
	if req.Date != nil {
		tx.Date = *req.Date
	}
	tx.Notes = req.Notes
	tx.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, tx); err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to update transaction", http.StatusInternalServerError)
	}
	return tx, nil
}

func (uc *financeUseCase) Delete(ctx context.Context, userID, id string) error {
	if _, err := uc.GetByID(ctx, userID, id); err != nil {
		return err
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		return apperror.New(apperror.CodeInternalError, "failed to delete transaction", http.StatusInternalServerError)
	}
	return nil
}

func (uc *financeUseCase) GetStats(ctx context.Context, userID string) (*in.FinanceStats, error) {
	opts := pagination.Params{Page: 1, Limit: 1000, Sort: "-date"}
	txs, _, err := uc.repo.FindByUserID(ctx, userID, opts)
	if err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to get finance stats", http.StatusInternalServerError)
	}

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	stats := &in.FinanceStats{ByCategory: make(map[string]float64)}

	for _, tx := range txs {
		stats.ByCategory[tx.Category] += tx.Amount
		if tx.Date.After(monthStart) {
			if tx.Type == domain.TransactionIncome {
				stats.ThisMonth.Income += tx.Amount
			} else {
				stats.ThisMonth.Expense += tx.Amount
			}
		}
	}
	stats.ThisMonth.Net = stats.ThisMonth.Income - stats.ThisMonth.Expense

	return stats, nil
}

func (uc *financeUseCase) GetCategories(ctx context.Context, userID string) ([]string, error) {
	opts := pagination.Params{Page: 1, Limit: 1000, Sort: "-created_at"}
	txs, _, err := uc.repo.FindByUserID(ctx, userID, opts)
	if err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "failed to get categories", http.StatusInternalServerError)
	}

	seen := make(map[string]struct{})
	var categories []string
	for _, tx := range txs {
		if _, ok := seen[tx.Category]; !ok {
			seen[tx.Category] = struct{}{}
			categories = append(categories, tx.Category)
		}
	}
	return categories, nil
}
