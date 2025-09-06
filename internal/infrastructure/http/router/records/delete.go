package records

import (
	"context"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) LedgerExpenseRecordDelete(ctx context.Context, params server.LedgerExpenseRecordDeleteParams) error {
	identity, err := a.GetIdentity(ctx)
	if err != nil {
		return err
	}

	apiReq := usercontroller.DeleteExpenseRecordRequest{
		ActorID:   identity.UserID,
		LedgerID:  domain.ConvertID(params.LedgerID),
		ExpenseID: domain.ConvertID(params.ExpenseID),
		RecordID:  domain.ConvertID(params.RecordID),
	}

	return a.Records().Delete(ctx, apiReq)
}
