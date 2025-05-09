package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/controllers"
	"github.com/sonalys/goshare/internal/domain"
)

func (a *API) LedgerExpenseCreate(ctx context.Context, req *handlers.Expense, params handlers.LedgerExpenseCreateParams) (r *handlers.LedgerExpenseCreateOK, _ error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	pendingRecords, err := convertUserBalances(req.Records)
	if err != nil {
		return nil, err
	}

	apiReq := controllers.CreateExpenseRequest{
		Actor:          identity.UserID,
		LedgerID:       domain.ConvertID(params.LedgerID),
		Name:           req.Name,
		ExpenseDate:    req.ExpenseDate,
		PendingRecords: pendingRecords,
	}

	switch resp, err := a.Ledgers.CreateExpense(ctx, apiReq); err {
	case nil:
		return &handlers.LedgerExpenseCreateOK{
			ID: resp.ID.UUID(),
		}, nil
	default:
		return nil, err
	}
}

func convertUserBalances(userBalances []handlers.ExpenseRecord) ([]domain.PendingRecord, error) {
	var errs domain.FormError

	balances := make([]domain.PendingRecord, 0, len(userBalances))
	for i, ub := range userBalances {
		recordType, err := domain.NewRecordType(string(ub.Type))
		if err != nil {
			errs = append(errs, domain.FieldError{
				Cause: err,
				Field: "records",
				Metadata: domain.FieldErrorMetadata{
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
