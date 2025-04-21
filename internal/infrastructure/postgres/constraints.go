package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	constraintParticipantUniqueEmail  = "participant_unique_email"
	constraintLedgerUniqueParticipant = "ledger_participant_unique"
	constraintLedgerParticipantsFK    = "ledger_participants_ledger_id_fkey"
	constraintLedgerRecordsUser       = "ledger_records_user_id_fkey"
)

func isConstraintError(err error) bool {
	if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
