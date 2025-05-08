package domain

type (
	Topic string

	Event[T any] struct {
		Topic Topic
		Data  T
	}

	GenericEvent interface {
		GetData() any
		GetTopic() Topic
	}
)

const (
	TopicLedgerCreated          Topic = "ledger.created"
	TopicLedgerExpenseCreated   Topic = "ledger.expense.created"
	TopicLedgerParticipantAdded Topic = "ledger.participant.added"
	TopicUserCreated            Topic = "user.created"
)

func (e Event[T]) GetTopic() Topic {
	return e.Topic
}

func (e Event[T]) GetData() any {
	return e.Data
}

var (
	_ GenericEvent = &Event[any]{}
)
