package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
)

func (a *API) LedgerExpenseRecordCreate(ctx context.Context, req *handlers.LedgerExpenseRecordCreateReq, params handlers.LedgerExpenseRecordCreateParams) (*handlers.Expense, error) {
	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	pendingRecords, err := convertUserBalances(req.Records)
	if err != nil {
		return nil, err
	}

	apiReq := usercontroller.CreateExpenseRecordRequest{
		Actor:          identity.UserID,
		LedgerID:       domain.ConvertID(params.LedgerID),
		ExpenseID:      domain.ConvertID(params.ExpenseID),
		PendingRecords: pendingRecords,
	}

	resp, err := a.UserController.Records().Create(ctx, apiReq)
	if err != nil {
		return nil, err
	}

	return &handlers.Expense{
		ID:          handlers.NewOptUUID(resp.ID.UUID()),
		Name:        resp.Name,
		ExpenseDate: resp.ExpenseDate,
		Records:     convertRecords(resp.Records),
	}, nil
}
