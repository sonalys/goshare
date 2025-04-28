package ledgers

import (
	"context"
	"time"

	"github.com/sonalys/goshare/internal/pkg/otel"
	"github.com/sonalys/goshare/internal/pkg/slog"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type (
	CreateRequest struct {
		UserID v1.ID
		Name   string
	}

	CreateResponse struct {
		ID v1.ID
	}
)

func (r CreateRequest) Validate() error {
	var errs v1.FormError

	if r.UserID.IsEmpty() {
		errs = append(errs, v1.NewRequiredFieldError("user_id"))
	}

	if r.Name == "" {
		errs = append(errs, v1.NewRequiredFieldError("name"))
	} else if nameLength := len(r.Name); nameLength < 3 || nameLength > 255 {
		errs = append(errs, v1.NewFieldLengthError("name", 3, 255))
	}

	return errs.Validate()
}

func (c *Controller) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ledgers.Create")
	defer span.End()

	ctx = slog.Context(ctx,
		slog.WithStringer("user_id", req.UserID),
	)

	if err := req.Validate(); err != nil {
		return nil, slog.ErrorReturn(ctx, "invalid request", err)
	}

	ledger := &v1.Ledger{
		ID:   v1.NewID(),
		Name: req.Name,
		Participants: []v1.LedgerParticipant{
			{
				ID:        v1.NewID(),
				UserID:    req.UserID,
				Balance:   0,
				CreatedAt: time.Now(),
				CreatedBy: req.UserID,
			},
		},
		CreatedAt: time.Now(),
		CreatedBy: req.UserID,
	}

	ctx = slog.Context(ctx,
		slog.WithStringer("ledger_id", ledger.ID),
	)

	err := c.ledgerRepository.Create(ctx, req.UserID, func(count int64) (*v1.Ledger, error) {
		if count+1 > v1.UserMaxLedgers {
			return nil, v1.ErrUserMaxLedgers
		}

		return ledger, nil
	})
	if err != nil {
		return nil, slog.ErrorReturn(ctx, "failed to create ledger", err)
	}

	slog.Info(ctx, "ledger created")

	resp := &CreateResponse{
		ID: ledger.ID,
	}

	return resp, nil
}
