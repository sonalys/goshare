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

func (s *Subscriber) handle(ctx context.Context, db Database, event Event) error {
	slog.Debug(ctx, "event created", slog.WithAny("event", event))
	return db.Transaction(ctx, func(db Database) error {
		for _, subscription := range s.subscriptions[event.GetTopic()] {
			if err := subscription(ctx, event, db); err != nil {
				return err
			}
		}
		return nil
	})
}
