package domain_test

import (
	"testing"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestParseID(t *testing.T) {
	t.Parallel()

	t.Run("pass/converts roundtrip", func(t *testing.T) {
		t.Parallel()
		id := domain.NewID()
		got, err := domain.ParseID(id.String())
		require.NoError(t, err)
		require.Equal(t, id, got)
	})

	t.Run("fail/invalid id", func(t *testing.T) {
		t.Parallel()
		got, err := domain.ParseID("a")
		require.Error(t, err)
		require.Empty(t, got)
	})
}
