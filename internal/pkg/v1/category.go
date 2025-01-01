package v1

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID
	LedgerID  uuid.UUID
	ParentID  uuid.UUID
	Name      string
	CreatedAt time.Time
	CreatedBy uuid.UUID
}
