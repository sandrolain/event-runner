package sources

import "time"

type ConsumerMessage struct {
	Type    string
	Topic   string
	Payload []byte
	Time    time.Time
}

type ConsumerCallback func(*ConsumerMessage, error) error

type Source interface {
	NewConsumer(ConsumerCallback)
}
