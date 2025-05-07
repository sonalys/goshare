package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sonalys/goshare/internal/domain"
)

func convertID(from domain.ID) pgtype.UUID {
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

func convertTime(from time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  from,
		Valid: true,
	}
}
