package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/sonalys/goshare/pkg/slog"
)

type (
	id = uuid.UUID
	ID struct{ id }
)

func NewID() ID {
	id, err := uuid.NewV7()
	if err != nil {
		slog.Panic(context.Background(), "could not generate uuid v7", err)
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
