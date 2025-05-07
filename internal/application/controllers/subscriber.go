package controllers

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
	"github.com/sonalys/goshare/internal/domain"
)

type (
	Event interface {
		GetTopic() domain.Topic
		GetData() any
	}

	Subscription func(ctx context.Context, event Event, uow Repositories) error

	Subscriber struct {
		db            Database
		subscriptions map[domain.Topic][]Subscription
	}
)

func newSubscriber(db Database) *Subscriber {
	return &Subscriber{
		db: db,
		subscriptions: map[domain.Topic][]Subscription{
			domain.TopicUserCreated: {onUserCreated},
		},
	}
}

func (s *Subscriber) handle(ctx context.Context, event Event) error {
	return s.db.Transaction(ctx, func(uow Repositories) error {
		for _, subscription := range s.subscriptions[event.GetTopic()] {
			if err := subscription(ctx, event, uow); err != nil {
				return err
			}
		}
		return nil
	})
}

func onUserCreated(ctx context.Context, event Event, uow Repositories) error {
	data, ok := event.GetData().(domain.UserCreated)
	if !ok {
		return fmt.Errorf("unexpected event type %T", event)
	}

	return uow.User().Create(ctx, &v1.User{
		ID:              v1.ConvertID(data.ID),
		FirstName:       data.FirstName,
		LastName:        data.LastName,
		Email:           data.Email,
		IsEmailVerified: false,
		PasswordHash:    data.HashedPassword,
		CreatedAt:       time.Now(),
	})
}
