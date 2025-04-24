package ledgers

import (
	"context"
	"log/slog"
	"time"

	"github.com/sonalys/goshare/internal/pkg/otel"
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

	if err := req.Validate(); err != nil {
		slog.ErrorContext(ctx, "invalid request", slog.Any("error", err))
		return nil, err
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

	err := c.ledgerRepository.Create(ctx, req.UserID, func(count int64) (*v1.Ledger, error) {
		if count+1 > v1.UserMaxLedgers {
			return nil, v1.ErrUserMaxLedgers
		}

		return ledger, nil
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to create ledger", slog.Any("error", err))
		return nil, err
	}

	slog.InfoContext(ctx, "ledger created", slog.String("ledger_id", ledger.ID.String()))

	resp := &CreateResponse{
		ID: ledger.ID,
	}

	return resp, nil
}
