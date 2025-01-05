package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (r *LedgerRepository) AppendRecord(ctx context.Context, ledgerID uuid.UUID, record *v1.LedgerRecord) error {
	return mapLedgerError(r.client.queries().AppendLedgerRecord(ctx, sqlc.AppendLedgerRecordParams{
		ID:          convertUUID(record.ID),
		LedgerID:    convertUUID(ledgerID),
		ExpenseID:   convertUUID(record.ExpenseID),
		UserID:      convertUUID(record.UserID),
		Description: record.Description,
		Amount:      record.Amount,
		CreatedAt:   convertTime(record.CreatedAt),
		CreatedBy:   convertUUID(record.CreatedBy),
	}))
}

func (r *LedgerRepository) GetRecords(ctx context.Context, ledgerID uuid.UUID) ([]v1.LedgerRecord, error) {
	records, err := r.client.queries().GetLedgerRecords(ctx, convertUUID(ledgerID))
	if err != nil {
		return nil, mapLedgerError(err)
	}
	result := make([]v1.LedgerRecord, 0, len(records))
	for _, record := range records {
		result = append(result, *newLedgerRecord(&record))
	}
	return result, nil
}

func newLedgerRecord(record *sqlc.LedgerRecord) *v1.LedgerRecord {
	return &v1.LedgerRecord{
		ID:          newUUID(record.ID),
		ExpenseID:   newUUID(record.ExpenseID),
		LedgerID:    newUUID(record.LedgerID),
		UserID:      newUUID(record.UserID),
		Amount:      record.Amount,
		CreatedAt:   record.CreatedAt.Time,
		CreatedBy:   newUUID(record.CreatedBy),
		Description: record.Description,
	}
}
