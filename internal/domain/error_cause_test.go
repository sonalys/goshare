package domain_test

import (
	"testing"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestCause_Is(t *testing.T) {
	t.Parallel()

	t.Run("pass/is", func(t *testing.T) {
		t.Parallel()
		//nolint
		require.ErrorIs(t, domain.ErrOverflow, domain.ErrOverflow)
	})

	t.Run("fail/isnt", func(t *testing.T) {
		t.Parallel()
		require.NotErrorIs(t, domain.ErrInvalid, domain.ErrOverflow)
	})
}
