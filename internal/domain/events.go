package domain

type Topic string

const (
	TopicLedgerCreated          Topic = "ledger.created"
	TopicLedgerExpenseCreated   Topic = "ledger.expense.created"
	TopicLedgerParticipantAdded Topic = "ledger.participant.added"
	TopicUserCreated            Topic = "user.created"
)

type Event[T any] struct {
	Topic Topic
	Data  T
}

func (e Event[T]) GetTopic() Topic {
	return e.Topic
}

func (e Event[T]) GetData() any {
	return e.Data
}
