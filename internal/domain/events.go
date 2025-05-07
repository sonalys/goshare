package domain

type Topic string

const (
	TopicUserCreated          Topic = "user.created"
	TopicLedgerCreated        Topic = "ledger.created"
	TopicLedgerExpenseCreated Topic = "ledger.expense.created"
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
