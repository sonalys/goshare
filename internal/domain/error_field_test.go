package domain_test

import (
	"testing"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestFieldErrorList_Error(t *testing.T) {
	t.Run("pass/more than one error", func(t *testing.T) {
		testData := domain.FieldErrorList{
			domain.FieldError{},
			domain.FieldError{},
		}

		msg := testData.Error()
		require.NotEmpty(t, msg)
	})
}
