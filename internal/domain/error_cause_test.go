package domain_test

import (
	"testing"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestCause_Is(t *testing.T) {
	t.Run("pass/is", func(t *testing.T) {
		require.ErrorIs(t, domain.CauseOverflow, domain.CauseOverflow)
	})

	t.Run("fail/isnt", func(t *testing.T) {
		require.NotErrorIs(t, domain.CauseInvalid, domain.CauseNotFound)
	})
}
