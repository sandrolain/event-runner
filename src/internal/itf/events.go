package itf

import (
	"time"

	"github.com/sandrolain/event-runner/src/config"
)

type EventConnection interface {
	NewInput(config.Input) (EventInput, error)
	NewOutput(config.Output) (EventOutput, error)
	NewCache(config.Cache) (EventCache, error)
	Close() error
}

// EventInput is a source of events that can be consumed by the EventRunner
//
// EventInput is an interface for sources of events that can be consumed
// by the EventRunner. This interface is implemented by sources such as
// Kafka or Redis.
//
// Receive is called by the EventRunner to receive events from the
// source. The implementation of this method should block until an event
// is available or an error occurs. When an event is available the
// implementation should call the `EventReceiver` function with the
// event. When an error occurs the implementation should return the
// error.
type EventInput interface {
	Receive() (chan EventMessage, error)
	Close() error
}

// EventMessage is an interface for messages that contain event data.
// The interface is implemented by event sources such as Kafka or Redis.
//
// The methods on EventMessage should be called by the EventRunner when it
// receives an event. The methods provide access to the details of the event
// such as the source of the event, the headers, and the data
type EventMessage interface {
	Time() (time.Time, error)
	Topic() (string, error)
	ReplyTo() (string, error)
	Metadata(string) ([]string, error)
	Data() ([]byte, error)
	DataString() (string, error)
	Ack() error
	Nak() error
}

type EventOutput interface {
	Ingest(chan RunnerResult) error
	Close() error
}

type EventCache interface {
	Get(key string) (any, error)
	Set(key string, data any) error
	Del(key string) error
	Close() error
}
