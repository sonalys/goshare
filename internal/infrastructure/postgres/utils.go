package postgres

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertTime(from time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  from,
		Valid: true,
	}
}
