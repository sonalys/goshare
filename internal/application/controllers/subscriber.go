package controllers

import (
	"context"

	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/domain"
)

type (
	Subscription func(ctx context.Context, event domain.GenericEvent, r Repositories) error

	Subscriber struct {
		subscriptions map[domain.Topic][]Subscription
	}
)

func NewSubscriber() *Subscriber {
	return &Subscriber{
		subscriptions: map[domain.Topic][]Subscription{},
	}
}

func convertEvents[T domain.GenericEvent](events []T) []domain.GenericEvent {
	out := make([]domain.GenericEvent, len(events))

	for i := range events {
		out[i] = events[i]
	}

	return out
}

func (s *Subscriber) Handle(ctx context.Context, db Database, events ...domain.GenericEvent) error {
	return db.Transaction(ctx, func(db Database) error {
		for _, event := range events {
			slog.Debug(ctx, "event created", slog.WithAny("event", event))
			for _, subscription := range s.subscriptions[event.GetTopic()] {
				if err := subscription(ctx, event, db); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
