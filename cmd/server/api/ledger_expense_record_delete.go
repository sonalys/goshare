package api

import (
	"context"

	"github.com/sonalys/goshare/cmd/server/handlers"
	"github.com/sonalys/goshare/internal/application/controllers"
	"github.com/sonalys/goshare/internal/domain"
)

func (a *API) LedgerExpenseRecordDelete(ctx context.Context, params handlers.LedgerExpenseRecordDeleteParams) error {
	identity, err := getIdentity(ctx)
	if err != nil {
		return err
	}

	apiReq := controllers.DeleteExpenseRecordRequest{
		ActorID:   identity.UserID,
		LedgerID:  domain.ConvertID(params.LedgerID),
		ExpenseID: domain.ConvertID(params.ExpenseID),
		RecordID:  domain.ConvertID(params.RecordID),
	}

	return a.Ledgers.DeleteExpenseRecord(ctx, apiReq)
}
