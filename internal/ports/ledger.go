package ports

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
)

type (
	LedgerQueries interface {
		// Get returns the ledger by id.
		// Returns application.ErrNotFound if it doesn't exist.
		Get(ctx context.Context, id domain.ID) (*domain.Ledger, error)
		// ListByUser returns all ledgers that the identity created or is a member.
		// Returns empty list if nothing is found.
		ListByUser(ctx context.Context, identity domain.ID) ([]domain.Ledger, error)
	}

	LedgerCommands interface {
		Create(ctx context.Context, ledger *domain.Ledger) error
		Update(ctx context.Context, ledger *domain.Ledger) error
	}

	LedgerRepository interface {
		LedgerQueries
		LedgerCommands
	}
)
