package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type (
	id = uuid.UUID
	ID struct{ id }
)

var EmptyID = ID{uuid.Nil}

func NewID() ID {
	id, err := uuid.NewV7()
	if err != nil {
		panic(fmt.Errorf("could not generate uuid v7: %w", err))
	}
	return ID{id: id}
}

func ParseID(from string) (ID, error) {
	id, err := uuid.Parse(from)
	if err != nil {
		return ID{uuid.Nil}, err
	}
	return ID{id}, nil
}

func (id *ID) UUID() uuid.UUID {
	return id.id
}

func (id *ID) IsEmpty() bool {
	return id.id == uuid.Nil
}

func ConvertID(id uuid.UUID) ID {
	return ID{id}
}
