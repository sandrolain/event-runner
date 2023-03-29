package kafka

import (
	"context"
	"fmt"

	"github.com/sandrolain/event-runner/src/sources"
	"github.com/segmentio/kafka-go"
)

const (
	EVENTS_TOPIC = "incoming"
	RESULT_TOPIC = "result"
)

func NewSource() sources.Source {
	return &KafkaBroker{}
}

type KafkaBroker struct {
	reader *kafka.Reader
}

func (k *KafkaBroker) NewConsumer(cb sources.ConsumerCallback) {
	// make a new reader that consumes from topic-A, partition 0, at offset 42
	k.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9093"},
		GroupID: "incomi",
		Topic:   EVENTS_TOPIC,
		// Partition: 0,
		//MinBytes:  10e2, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer k.reader.Close()

	fmt.Printf("Kafka source on %v\n", k.reader.Config().Brokers)

	for {
		m, err := k.reader.ReadMessage(context.Background())
		// m, err := k.reader.FetchMessage(context.Background())
		if err != nil {
			cb(nil, err)
			continue
		}
		err = cb(&sources.ConsumerMessage{
			Type:    "kafka",
			Topic:   m.Topic,
			Payload: m.Value,
			Time:    m.Time,
		}, nil)
		// if err == nil {
		// 	k.reader.CommitMessages(context.Background(), m)
		// }
	}
}
