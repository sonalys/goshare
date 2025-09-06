package expenses

import (
	"context"
	"time"

	"github.com/sonalys/goshare/internal/application/controllers/usercontroller"
	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/http/mappers"
	"github.com/sonalys/goshare/internal/infrastructure/http/server"
)

func (a *Router) LedgerExpenseList(ctx context.Context, params server.LedgerExpenseListParams) (*server.LedgerExpenseListOK, error) {
	identity, err := a.GetIdentity(ctx)
	if err != nil {
		return nil, err
	}

	result, err := a.Expenses().List(ctx, usercontroller.ListExpensesRequest{
		ActorID:  identity.UserID,
		LedgerID: domain.ConvertID(params.LedgerID),
		Limit:    params.Limit.Or(10),
		Cursor:   params.Cursor.Or(time.Now()),
	})
	if err != nil {
		return nil, err
	}

	var cursor server.OptDateTime

	if result.Cursor != nil {
		cursor = server.NewOptDateTime(*result.Cursor)
	}

	return &server.LedgerExpenseListOK{
		Expenses: mappers.LedgerExpenseSummaryToExpenseSummary(result.Expenses),
		Cursor:   cursor,
	}, nil
}
