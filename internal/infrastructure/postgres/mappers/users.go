package mappers

import (
	domain "github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
)

func NewUser(user sqlc.User) *domain.User {
	return &domain.User{
		ID:              user.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		IsEmailVerified: false,
		PasswordHash:    user.PasswordHash,
		CreatedAt:       user.CreatedAt.Time,
		LedgersCount:    user.LedgerCount,
	}
}

func NewUsers(from []sqlc.User) []domain.User {
	to := make([]domain.User, 0, len(from))

	for i := range from {
		to = append(to, *NewUser(from[i]))
	}

	return to
}
