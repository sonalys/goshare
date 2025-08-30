package v1

import (
	"time"

	"github.com/sonalys/goshare/internal/domain"
)

type (
	LedgerExpenseSummary struct {
		ID          domain.ID
		Amount      int32
		Name        string
		ExpenseDate time.Time
		CreatedAt   time.Time
		CreatedBy   domain.ID
		UpdatedAt   time.Time
		UpdatedBy   domain.ID
	}
)
