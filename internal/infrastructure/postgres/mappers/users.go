package mappers

import (
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
)

func NewUser(user sqlc.User) *v1.User {
	return &v1.User{
		ID:              newUUID(user.ID),
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		IsEmailVerified: false,
		PasswordHash:    user.PasswordHash,
		CreatedAt:       user.CreatedAt.Time,
	}
}

func NewUsers(from []sqlc.User) []v1.User {
	to := make([]v1.User, 0, len(from))

	for i := range from {
		to = append(to, *NewUser(from[i]))
	}

	return to
}
