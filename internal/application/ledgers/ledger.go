package ledgers

import (
	"context"

	"github.com/sonalys/goshare/internal/pkg/otel"
	"github.com/sonalys/goshare/internal/pkg/slog"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (c *Controller) GetByUser(ctx context.Context, userID v1.ID) ([]v1.Ledger, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.ListByUser")
	defer span.End()

	ledgers, err := c.ledgerRepository.GetByUser(ctx, userID)
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "failed to list ledgers", err)
	}

	return ledgers, nil
}
