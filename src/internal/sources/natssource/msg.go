package natssource

import (
	"time"

	"github.com/nats-io/nats.go"
)

type NatsEventMessage struct {
	time time.Time
	msg  *nats.Msg
}

func (m NatsEventMessage) Topic() (string, error) {
	return m.msg.Subject, nil
}

func (m NatsEventMessage) ReplyTo() (string, error) {
	return m.msg.Reply, nil
}

func (m NatsEventMessage) Metadata(key string) (res []string, err error) {
	value := m.msg.Header.Get(key)
	return []string{value}, nil
}

func (m NatsEventMessage) Data() ([]byte, error) {
	return m.msg.Data, nil
}

func (m NatsEventMessage) Time() (time.Time, error) {
	return m.time, nil
}

func (m NatsEventMessage) Ack() error {
	return m.msg.Ack()
}

func (m NatsEventMessage) Nak() error {
	return m.msg.Nak()
}
