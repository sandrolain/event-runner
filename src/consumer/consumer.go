package consumer

import (
	"fmt"

	"github.com/go-redis/redis"
)

func RedisConnection(fn func(*redis.Message)) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "development.password", // no password set
		DB:       0,                      // use default DB
	})

	pubsub := rdb.Subscribe("events")
	// Get the Channel to use
	channel := pubsub.Channel()

	go redisListener(channel, fn)

	fmt.Println("Redis listener setup")

	return rdb
}

func redisListener(channel <-chan *redis.Message, fn func(*redis.Message)) {
	for msg := range channel {
		fn(msg)
	}
}
