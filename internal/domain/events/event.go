package events

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type (
	Type string

	Event interface {
		ID() string
		Timestamp() time.Time
		Type() Type
		Version() int64
		SetVersion(int64)
		Content(dst any) error
	}

	event struct {
		id            string
		t             Type
		timestamp     time.Time
		version       int64
		contentType   string
		content       []byte
		previousEvent *string
	}

	contextKey string
)

const (
	currentEventKey = contextKey("currentEvent")
)

func New() Event {
	return &event{}
}

func Load(
	ctx context.Context,
	id string,
) (context.Context, Event) {
	return setContext(ctx, id), &event{
		id:            id,
		previousEvent: getContextEvent(ctx),
	}
}

func setContext(ctx context.Context, eventID string) context.Context {
	return context.WithValue(ctx, currentEventKey, eventID)
}

func getContextEvent(ctx context.Context) *string {
	eventID, ok := ctx.Value(currentEventKey).(string)
	if ok {
		return &eventID
	}

	return nil
}

func (e *event) Content(dst any) error {
	switch e.contentType {
	case "application/json":
		return json.Unmarshal(e.content, dst)
	default:
		return errors.New("invalid content-type")
	}
}

func (e *event) ID() string {
	return e.id
}

func (e *event) Timestamp() time.Time {
	return e.timestamp
}

func (e *event) SetVersion(v int64) {
	e.version = v
}

func (e *event) Version() int64 {
	return e.version
}

func (e *event) Type() Type {
	return e.t
}
