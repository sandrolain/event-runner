package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/sandrolain/event-runner/src/sources"
)

const (
	EVENTS_CHANNEL = "incoming"
	RESULT_CHANNEL = "result"
)

func NewSource() sources.Source {
	return &RedisBroker{}
}

type RedisBroker struct {
	client *redis.Client
}

func (r *RedisBroker) NewConsumer(cb sources.ConsumerCallback) {
	r.client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "development.password", // no password set
		DB:       0,                      // use default DB
	})
	defer r.client.Close()

	fmt.Printf("Redis source on %s\n", r.client.Options().Addr)

	pubsub := r.client.Subscribe(EVENTS_CHANNEL)
	// Get the Channel to use
	channel := pubsub.Channel()

	for msg := range channel {
		cb(&sources.ConsumerMessage{
			Time:    time.Now(),
			Topic:   msg.Channel,
			Payload: []byte(msg.Payload),
		}, nil)
	}
}
