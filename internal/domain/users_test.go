package domain_test

import (
	"strings"
	"testing"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Parallel()

	factory := func(hooks ...func(req *domain.NewUserRequest)) domain.NewUserRequest {
		req := domain.NewUserRequest{
			FirstName: "First",
			LastName:  "Last",
			Email:     "email@domain.com",
			Password:  "password",
		}

		for _, hook := range hooks {
			hook(&req)
		}

		return req
	}

	t.Run("pass", func(t *testing.T) {
		t.Parallel()
		req := factory()

		got, err := domain.NewUser(req)

		require.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, req.FirstName, got.FirstName)
		assert.Equal(t, req.LastName, got.LastName)
		assert.Equal(t, req.Email, got.Email)
		assert.NotEmpty(t, got.PasswordHash)
		assert.False(t, got.IsEmailVerified)
		assert.Zero(t, got.LedgersCount)
		assert.NotZero(t, got.CreatedAt)
	})

	t.Run("fail/empty first name", func(t *testing.T) {
		t.Parallel()
		req := factory(func(req *domain.NewUserRequest) {
			req.FirstName = ""
		})

		got, err := domain.NewUser(req)

		assert.Nil(t, got)
		var targetErr domain.FieldError
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, "firstName", targetErr.Field)
		assert.Equal(t, domain.ErrCauseRequired, targetErr.Cause)
	})

	t.Run("fail/empty last name", func(t *testing.T) {
		t.Parallel()
		req := factory(func(req *domain.NewUserRequest) {
			req.LastName = ""
		})

		got, err := domain.NewUser(req)

		assert.Nil(t, got)
		var targetErr domain.FieldError
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, "lastName", targetErr.Field)
		assert.Equal(t, domain.ErrCauseRequired, targetErr.Cause)
	})

	t.Run("fail/short password", func(t *testing.T) {
		t.Parallel()
		req := factory(func(req *domain.NewUserRequest) {
			req.Password = ""
		})

		got, err := domain.NewUser(req)

		assert.Nil(t, got)
		var targetErr domain.FieldError
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, "password", targetErr.Field)
		assert.Equal(t, &domain.ValueLengthError{Min: 8, Max: 72}, targetErr.Cause)
	})

	t.Run("fail/long password", func(t *testing.T) {
		t.Parallel()
		req := factory(func(req *domain.NewUserRequest) {
			req.Password = strings.Repeat("a", 73)
		})

		got, err := domain.NewUser(req)

		assert.Nil(t, got)
		var targetErr domain.FieldError
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, "password", targetErr.Field)
		assert.Equal(t, &domain.ValueLengthError{Min: 8, Max: 72}, targetErr.Cause)
	})

	t.Run("fail/empty email", func(t *testing.T) {
		t.Parallel()
		req := factory(func(req *domain.NewUserRequest) {
			req.Email = ""
		})

		got, err := domain.NewUser(req)

		assert.Nil(t, got)
		var targetErr domain.FieldError
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, "email", targetErr.Field)
		assert.Equal(t, domain.ErrCauseInvalid, targetErr.Cause)
	})

	t.Run("fail/invalid email", func(t *testing.T) {
		t.Parallel()
		req := factory(func(req *domain.NewUserRequest) {
			req.Email = "invalid@"
		})

		got, err := domain.NewUser(req)

		assert.Nil(t, got)
		var targetErr domain.FieldError
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, "email", targetErr.Field)
		assert.Equal(t, domain.ErrCauseInvalid, targetErr.Cause)
	})
}

func TestUser_CreateLedger(t *testing.T) {
	t.Parallel()
	t.Run("pass", func(t *testing.T) {
		user := domain.User{
			ID: domain.NewID(),
		}

		ledger, err := user.CreateLedger("name")

		require.NoError(t, err)

		assert.Equal(t, ledger.Name, "name")
		assert.NotZero(t, ledger.ID)
		assert.NotZero(t, ledger.CreatedAt)
		assert.Equal(t, user.ID, ledger.CreatedBy)

		require.Len(t, ledger.Members, 1)

		participant := ledger.Members[0]

		assert.NotZero(t, participant.ID)
		assert.Equal(t, user.ID, participant.Identity)
		assert.Zero(t, participant.Balance)
		assert.NotZero(t, participant.CreatedAt)
		assert.Equal(t, user.ID, participant.CreatedBy)
	})

	t.Run("fail/short name", func(t *testing.T) {
		t.Parallel()
		user := domain.User{
			ID: domain.NewID(),
		}

		ledger, err := user.CreateLedger("")

		assert.Nil(t, ledger)
		var targetErr domain.FieldError
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, "name", targetErr.Field)
		assert.Equal(t, &domain.ValueLengthError{Min: 3, Max: 255}, targetErr.Cause)
	})

	t.Run("fail/long name", func(t *testing.T) {
		t.Parallel()
		user := domain.User{
			ID: domain.NewID(),
		}

		ledger, err := user.CreateLedger(strings.Repeat("a", 256))

		assert.Nil(t, ledger)
		var targetErr domain.FieldError
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, "name", targetErr.Field)
		assert.Equal(t, &domain.ValueLengthError{Min: 3, Max: 255}, targetErr.Cause)
	})

	t.Run("fail/user max ledgers", func(t *testing.T) {
		t.Parallel()
		user := domain.User{
			ID:           domain.NewID(),
			LedgersCount: domain.UserMaxLedgers,
		}

		ledger, err := user.CreateLedger("name")

		assert.Nil(t, ledger)
		var targetErr *domain.ErrUserMaxLedgers
		require.ErrorAs(t, err, &targetErr)
		assert.Equal(t, user.ID, targetErr.UserID)
		assert.Equal(t, domain.UserMaxLedgers, targetErr.MaxLedgers)
	})
}
