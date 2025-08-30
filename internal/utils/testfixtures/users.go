package testfixtures

import (
	"testing"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/stretchr/testify/require"
)

func User(t *testing.T) *domain.User {
	user, err := domain.NewUser(domain.NewUserRequest{
		FirstName: domain.NewID().String(),
		LastName:  domain.NewID().String(),
		Email:     domain.NewID().String() + "@example.com",
		Password:  domain.NewID().String(),
	})
	require.NoError(t, err)

	return user
}

func Ledger(t *testing.T, creator *domain.User) *domain.Ledger {
	ledger, err := creator.CreateLedger(domain.NewID().String())
	require.NoError(t, err)
	return ledger
}
