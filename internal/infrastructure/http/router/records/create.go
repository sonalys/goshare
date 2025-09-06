package records

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/mappers"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) LedgerExpenseRecordCreate(ctx context.Context, req *server.LedgerExpenseRecordCreateReq, params server.LedgerExpenseRecordCreateParams) (*server.Expense, error) {
	identity, err := a.GetIdentity(ctx)
	if err != nil {
		return nil, err
	}

	pendingRecords, err := mappers.ExpenseRecordToPendingRecord(req.Records)
	if err != nil {
		return nil, err
	}

	apiReq := usercontroller.CreateExpenseRecordRequest{
		ActorID:        identity.UserID,
		LedgerID:       domain.ConvertID(params.LedgerID),
		ExpenseID:      domain.ConvertID(params.ExpenseID),
		PendingRecords: pendingRecords,
	}

	resp, err := a.Records().Create(ctx, apiReq)
	if err != nil {
		return nil, err
	}

	return &server.Expense{
		ID:          server.NewOptUUID(resp.ID.UUID()),
		Name:        resp.Name,
		ExpenseDate: resp.ExpenseDate,
		Records:     mappers.RecordToExpenseRecord(resp.Records),
	}, nil
}
