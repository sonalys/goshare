package mappers

import (
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
)

func NewLedgerParticipant(user *sqlc.LedgerParticipant) *v1.LedgerParticipant {
	return &v1.LedgerParticipant{
		ID:        newUUID(user.ID),
		UserID:    newUUID(user.UserID),
		Balance:   user.Balance,
		CreatedAt: user.CreatedAt.Time,
		CreatedBy: newUUID(user.CreatedBy),
	}
}
