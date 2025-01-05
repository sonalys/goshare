package v1

import (
	"time"
)

type Category struct {
	ID        ID
	LedgerID  ID
	ParentID  ID
	Name      string
	CreatedAt time.Time
	CreatedBy ID
}
