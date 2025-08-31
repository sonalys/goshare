package ledgers

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) LedgerExpenseCreate(ctx context.Context, req *server.Expense, params server.LedgerExpenseCreateParams) (r *server.LedgerExpenseCreateOK, _ error) {
	identity, err := a.GetIdentity(ctx)
	if err != nil {
		return nil, err
	}

	pendingRecords, err := convertUserBalances(req.Records)
	if err != nil {
		return nil, err
	}

	apiReq := usercontroller.CreateExpenseRequest{
		ActorID:        identity.UserID,
		LedgerID:       domain.ConvertID(params.LedgerID),
		Name:           req.Name,
		ExpenseDate:    req.ExpenseDate,
		PendingRecords: pendingRecords,
	}

	switch resp, err := a.Expenses().Create(ctx, apiReq); err {
	case nil:
		return &server.LedgerExpenseCreateOK{
			ID: resp.ID.UUID(),
		}, nil
	default:
		return nil, err
	}
}

func convertUserBalances(userBalances []server.ExpenseRecord) ([]domain.PendingRecord, error) {
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
