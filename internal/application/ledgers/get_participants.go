package ledgers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sonalys/goshare/internal/pkg/otel"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (c *Controller) GetParticipants(ctx context.Context, ledgerID v1.ID) ([]v1.LedgerParticipant, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.GetBalances")
	defer span.End()

	participants, err := c.ledgerRepository.GetParticipants(ctx, ledgerID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get ledger participants balances", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get ledger participants balances: %w", err)
	}

	slog.InfoContext(ctx, "ledger participants balances retrieved")

	return participants, nil
}
