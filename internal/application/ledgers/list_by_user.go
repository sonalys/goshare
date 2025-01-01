package ledgers

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/pkg/otel"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
	"go.opentelemetry.io/otel/codes"
)

func (c *Controller) ListByUser(ctx context.Context, userID uuid.UUID) ([]v1.Ledger, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.ListByUser")
	defer span.End()

	ledgers, err := c.repository.GetByUser(ctx, userID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed to list ledgers", slog.Any("error", err))
		return nil, err
	}

	return ledgers, nil
}
