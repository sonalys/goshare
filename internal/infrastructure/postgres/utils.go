package postgres

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertTime(from time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  from,
		Valid: true,
	}
}
