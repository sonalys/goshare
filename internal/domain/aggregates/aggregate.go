package aggregates

import (
	"errors"
	"fmt"

	"github.com/sonalys/goshare/internal/domain/events"
)

type (
	When interface {
		When(event events.Event) error
	}

	Load interface {
		Load(evts ...events.Event) error
	}

	Type string

	Aggregate interface {
		When
		Load
	}

	Base struct {
		ID               string
		Version          int64
		UncommitedEvents []events.Event
		Type             Type
		Handler          When
	}
)

func (a *Base) Load(evts ...events.Event) error {
	for _, event := range evts {
		if event.ID() != a.ID {
			return errors.New("invalid aggregate")
		}

		if err := a.Handler.When(event); err != nil {
			return fmt.Errorf("loading applying event: %w", err)
		}

		a.Version++
	}

	return nil
}

func (a *Base) Apply(event events.Event) error {
	if err := a.Load(event); err != nil {
		return err
	}

	event.SetVersion(a.Version)
	a.UncommitedEvents = append(a.UncommitedEvents, event)

	return nil
}
