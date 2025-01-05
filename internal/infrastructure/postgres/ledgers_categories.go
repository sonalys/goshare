package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (r *LedgerRepository) CreateCategory(ctx context.Context, ledgerID uuid.UUID, category *v1.Category) error {
	return mapLedgerError(r.client.queries().CreateCategory(ctx, sqlc.CreateCategoryParams{
		ID:        convertUUID(category.ID),
		LedgerID:  convertUUID(ledgerID),
		ParentID:  convertUUID(category.ParentID),
		Name:      category.Name,
		CreatedAt: convertTime(category.CreatedAt),
		CreatedBy: convertUUID(category.CreatedBy),
	}))
}

func (r *LedgerRepository) GetCategories(ctx context.Context, ledgerID uuid.UUID) ([]v1.Category, error) {
	categories, err := r.client.queries().GetLedgerCategories(ctx, convertUUID(ledgerID))
	if err != nil {
		return nil, mapLedgerError(err)
	}
	result := make([]v1.Category, 0, len(categories))
	for _, category := range categories {
		result = append(result, *newCategory(&category))
	}
	return result, nil
}

func newCategory(category *sqlc.Category) *v1.Category {
	return &v1.Category{
		ID:        newUUID(category.ID),
		LedgerID:  newUUID(category.LedgerID),
		ParentID:  newUUID(category.ParentID),
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Time,
		CreatedBy: newUUID(category.CreatedBy),
	}
}
