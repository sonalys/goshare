package domain_test

import (
	"testing"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/stretchr/testify/assert"
)

//  Replace "your_module_name" with your actual module name

func TestNewRecordType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    domain.RecordType
		wantErr bool
	}{
		{
			name:    "valid debt",
			input:   "debt",
			want:    domain.RecordTypeDebt,
			wantErr: false,
		},
		{
			name:    "valid settlement",
			input:   "settlement",
			want:    domain.RecordTypeSettlement,
			wantErr: false,
		},
		{
			name:    "invalid type",
			input:   "invalid",
			want:    domain.RecordTypeUnknown,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    domain.RecordTypeUnknown,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := domain.NewRecordType(tt.input)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestRecordType_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		r    domain.RecordType
		want string
	}{
		{
			name: "debt",
			r:    domain.RecordTypeDebt,
			want: "debt",
		},
		{
			name: "settlement",
			r:    domain.RecordTypeSettlement,
			want: "settlement",
		},
		{
			name: "unknown",
			r:    domain.RecordTypeUnknown,
			want: "unknown",
		},
		{
			name: "out of bounds",
			r:    domain.RecordType(99),
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.r.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRecordType_IsValid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		r    domain.RecordType
		want bool
	}{
		{
			name: "valid debt",
			r:    domain.RecordTypeDebt,
			want: true,
		},
		{
			name: "valid settlement",
			r:    domain.RecordTypeSettlement,
			want: true,
		},
		{
			name: "invalid unknown",
			r:    domain.RecordTypeUnknown,
			want: false,
		},
		{
			name: "invalid out of bounds",
			r:    domain.RecordType(99), // Assuming 99 is out of bounds
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.r.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}
