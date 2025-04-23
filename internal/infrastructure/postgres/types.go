package postgres

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func convertUUID(from v1.ID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: from.UUID(),
		Valid: true,
	}
}

func newUUID(from pgtype.UUID) v1.ID {
	return v1.ConvertID(from.Bytes)
}

func convertTime(from time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  from,
		Valid: true,
	}
}
