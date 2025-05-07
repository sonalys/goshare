package mappers

import (
	"github.com/jackc/pgx/v5/pgtype"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
)

func newUUID(from pgtype.UUID) v1.ID {
	return v1.ConvertID(from.Bytes)
}
