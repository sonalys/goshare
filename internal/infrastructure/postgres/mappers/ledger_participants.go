package mappers

import (
	domain "github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
)

func NewLedgerParticipant(user *sqlc.LedgerParticipant) *domain.LedgerMember {
	return &domain.LedgerMember{
		ID:        newUUID(user.ID),
		Identity:  newUUID(user.UserID),
		Balance:   user.Balance,
		CreatedAt: user.CreatedAt.Time,
		CreatedBy: newUUID(user.CreatedBy),
	}
}
