package usercontroller

import (
	"context"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/ports"
	"go.opentelemetry.io/otel/trace"
)

type (
	LedgerController interface {
		Create(ctx context.Context, req CreateLedgerRequest) (resp *CreateLedgerResponse, err error)
		Get(ctx context.Context, req GetLedgerRequest) (*domain.Ledger, error)
		ListByUser(ctx context.Context, actorID domain.ID) ([]domain.Ledger, error)
		MembersAdd(ctx context.Context, req AddMembersRequest) error
	}

	ledgerController struct {
		db     ports.LocalDatabase
		tracer trace.Tracer
	}
)
