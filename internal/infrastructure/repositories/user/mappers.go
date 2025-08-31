package user

import (
	domain "github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

func toUser(user sqlcgen.User) *domain.User {
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

func toUsers(from []sqlcgen.User) []domain.User {
	to := make([]domain.User, 0, len(from))

	for i := range from {
		to = append(to, *toUser(from[i]))
	}

	return to
}
