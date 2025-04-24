package ledgers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sonalys/goshare/internal/pkg/otel"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (c *Controller) GetByUser(ctx context.Context, userID v1.ID) ([]v1.Ledger, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.ListByUser")
	defer span.End()

	ledgers, err := c.ledgerRepository.GetByUser(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list ledgers", slog.Any("error", err))
		return nil, fmt.Errorf("failed to list ledgers: %w", err)
	}

	return ledgers, nil
}
