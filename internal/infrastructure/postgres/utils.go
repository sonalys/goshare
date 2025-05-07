package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
)

func convertID(from v1.ID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: from.UUID(),
		Valid: true,
	}
}

func convertUUID(from uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: from,
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
