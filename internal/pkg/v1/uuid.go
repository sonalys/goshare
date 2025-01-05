package v1

import (
	"fmt"

	"github.com/google/uuid"
)

type privateUUID = uuid.UUID

type ID struct{ privateUUID }

func NewID() ID {
	id, err := uuid.NewV7()
	if err != nil {
		panic(fmt.Errorf("could not generate uuid v7: %w", err))
	}
	return ID{privateUUID: id}
}

func ConvertID(id uuid.UUID) ID {
	return ID{id}
}

func ParseID(from string) (ID, error) {
	id, err := uuid.Parse(from)
	if err != nil {
		return ID{uuid.Nil}, err
	}
	return ID{id}, nil
}

func (id *ID) UUID() uuid.UUID {
	return id.privateUUID
}

func (id *ID) IsEmpty() bool {
	return id.privateUUID == uuid.Nil
}
