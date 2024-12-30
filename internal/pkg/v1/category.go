package v1

import "github.com/google/uuid"

type Category struct {
	ID        uuid.UUID
	ParentID  uuid.UUID
	Name      string
	CreatedAt string
	CreatedBy uuid.UUID
}
