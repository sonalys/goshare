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

func (c *Controller) GetByUser(ctx context.Context, userID uuid.UUID) ([]v1.Ledger, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.ListByUser")
	defer span.End()

	ledgers, err := c.ledgerRepository.GetByUser(ctx, userID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed to list ledgers", slog.Any("error", err))
		return nil, fmt.Errorf("failed to list ledgers: %w", err)
	}

	return ledgers, nil
}
