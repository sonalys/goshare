package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	constraintMemberUniqueEmail  = "member_unique_email"
	constraintLedgerUniqueMember = "ledger_member_unique"
	constraintLedgerMembersFK    = "ledger_members_ledger_id_fkey"
	constraintLedgerRecordsUser  = "ledger_records_user_id_fkey"
)

func isConstraintError(err error) bool {
	if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}

func isViolatingConstraint(err error, constraintName string) bool {
	if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
		return pgErr.ConstraintName == constraintName
	}
	return false
}
