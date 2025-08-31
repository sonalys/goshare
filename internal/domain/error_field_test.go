package domain_test

import (
	"testing"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestFieldErrorList_Error(t *testing.T) {
	t.Parallel()

	t.Run("pass/more than one error", func(t *testing.T) {
		t.Parallel()
		testData := domain.FieldErrors{
			domain.FieldError{},
			domain.FieldError{},
		}

		msg := testData.Error()
		require.NotEmpty(t, msg)
	})
}
