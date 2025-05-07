package controllers

import (
	"context"

	"github.com/sonalys/goshare/internal/application/pkg/slog"
	"github.com/sonalys/goshare/internal/domain"
)

type (
	Event interface {
		GetTopic() domain.Topic
		GetData() any
	}

	Subscription func(ctx context.Context, event Event, r Repositories) error

	Subscriber struct {
		subscriptions map[domain.Topic][]Subscription
	}
)

func newSubscriber() *Subscriber {
	return &Subscriber{
		subscriptions: map[domain.Topic][]Subscription{},
	}
}

func convertEvents[T Event](events []T) []Event {
	out := make([]Event, 0, len(events))

	for i := range events {
		out = append(out, events[i])
	}

	return out
}

func (s *Subscriber) handle(ctx context.Context, db Database, events ...Event) error {
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
