package expenses

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/mappers"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) LedgerExpenseCreate(ctx context.Context, req *server.Expense, params server.LedgerExpenseCreateParams) (r *server.LedgerExpenseCreateOK, _ error) {
	identity, err := a.securityHandler.GetIdentity(ctx)
	if err != nil {
		return nil, err
	}

	pendingRecords, err := mappers.ExpenseRecordToPendingRecord(req.Records)
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

	switch resp, err := a.controller.Expenses().Create(ctx, apiReq); err {
	case nil:
		return &server.LedgerExpenseCreateOK{
			ID: resp.ID.UUID(),
		}, nil
	default:
		return nil, err
	}
}
