package postgres

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func convertTime(from time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  from,
		Valid: true,
	}
}
