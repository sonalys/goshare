package mappers

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sonalys/goshare/internal/domain"
)

func newUUID(from pgtype.UUID) domain.ID {
	return domain.ConvertID(from.Bytes)
}
