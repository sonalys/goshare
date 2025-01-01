package ledgers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/pkg/otel"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
	"go.opentelemetry.io/otel/codes"
)

func (c *Controller) GetBalances(ctx context.Context, ledgerID uuid.UUID) ([]v1.LedgerParticipantBalance, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.GetBalances")
	defer span.End()

	balances, err := c.ledgerRepository.GetLedgerBalance(ctx, ledgerID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed to get ledger participants balances", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get ledger participants balances: %w", err)
	}

	span.SetStatus(codes.Ok, "")
	slog.InfoContext(ctx, "ledger participants balances retrieved")

	return balances, nil
}
