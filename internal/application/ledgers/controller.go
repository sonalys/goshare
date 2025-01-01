package ledgers

import (
	"context"

	"github.com/google/uuid"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

type (
	Repository interface {
		Create(ctx context.Context, ledger *v1.Ledger) error
		GetByUser(ctx context.Context, userID uuid.UUID) ([]v1.Ledger, error)
		GetLedgerParticipantsBalances(ctx context.Context, ledgerID uuid.UUID) ([]v1.LedgerParticipantBalance, error)
	}

	Controller struct {
		repository Repository
	}
)

func NewController(
	repository Repository,
) *Controller {
	return &Controller{
		repository: repository,
	}
}
