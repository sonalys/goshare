package mappers

import (
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func ExpenseRecordToPendingRecord(userBalances []server.ExpenseRecord) ([]domain.PendingRecord, error) {
	var errs domain.Form

	balances := make([]domain.PendingRecord, 0, len(userBalances))
	for i, ub := range userBalances {
		recordType, err := domain.NewRecordType(string(ub.Type))
		if err != nil {
			errs.Append(domain.FieldError{
				Cause: err,
				Field: "records",
				Metadata: &domain.FieldErrorMetadata{
					Index: i,
				},
			})
		}

		balances = append(balances, domain.PendingRecord{
			Type:   recordType,
			Amount: ub.Amount,
			From:   domain.ConvertID(ub.FromUserID),
			To:     domain.ConvertID(ub.ToUserID),
		})
	}

	if err := errs.Close(); err != nil {
		return nil, err
	}

	return balances, nil
}

func RecordToExpenseRecord(records map[domain.ID]*domain.Record) []server.ExpenseRecord {
	result := make([]server.ExpenseRecord, 0, len(records))
	for id, record := range records {
		result = append(result, server.ExpenseRecord{
			ID:         server.NewOptUUID(id.UUID()),
			Type:       server.ExpenseRecordType(record.Type.String()),
			FromUserID: record.From.UUID(),
			ToUserID:   record.To.UUID(),
			Amount:     record.Amount,
		})
	}

	return result
}
