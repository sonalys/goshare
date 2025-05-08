package domain

type (
	Topic string

	event[T any] struct {
		topic Topic
		data  T
	}

	Event interface {
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

func (e event[T]) GetTopic() Topic {
	return e.topic
}

func (e event[T]) GetData() any {
	return e.data
}

var (
	_ Event = &event[any]{}
)
