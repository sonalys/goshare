package ledgers

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/internal/pkg/otel"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
	"go.opentelemetry.io/otel/codes"
)

type (
	CreateRequest struct {
		UserID uuid.UUID
		Name   string
	}

	CreateResponse struct {
		ID uuid.UUID
	}
)

func (r CreateRequest) Validate() error {
	var errs v1.FormError

	if r.UserID == uuid.Nil {
		errs.Fields = append(errs.Fields, v1.NewRequiredFieldError("user_id"))
	}

	if r.Name == "" {
		errs.Fields = append(errs.Fields, v1.NewRequiredFieldError("name"))
	} else if nameLength := len(r.Name); nameLength < 3 || nameLength > 255 {
		errs.Fields = append(errs.Fields, v1.NewFieldLengthError("name", 3, 255))
	}

	return errs.Validate()
}

func (c *Controller) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.Create")
	defer span.End()

	if err := req.Validate(); err != nil {
		span.SetStatus(codes.Error, "invalid request")
		slog.ErrorContext(ctx, "invalid request", slog.Any("error", err))
		return nil, err
	}

	ledger := &v1.Ledger{
		ID:        uuid.New(),
		Name:      req.Name,
		CreatedAt: time.Now(),
		CreatedBy: req.UserID,
	}

	if err := c.repository.Create(ctx, ledger); err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed to create ledger", slog.Any("error", err))
		return nil, err
	}

	slog.InfoContext(ctx, "ledger created", slog.String("ledger", ledger.ID.String()))

	return &CreateResponse{
		ID: ledger.ID,
	}, nil
}