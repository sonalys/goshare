package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func convertUUID(from uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: from,
		Valid: true,
	}
}

func convertUUIDPtr(from *uuid.UUID) pgtype.UUID {
	if from == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{
		Bytes: *from,
		Valid: true,
	}
}

func newUUID(from pgtype.UUID) uuid.UUID {
	return uuid.UUID(from.Bytes)
}

func newUUIDPtr(from pgtype.UUID) *uuid.UUID {
	if !from.Valid {
		return nil
	}
	uuid := newUUID(from)
	return &uuid
}

func convertTime(from time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  from,
		Valid: true,
	}
}
