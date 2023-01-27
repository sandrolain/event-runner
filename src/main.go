package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sandrolain/event-runner/src/consumer"
	"github.com/sandrolain/event-runner/src/runner"
	"rogchap.com/v8go"
)

func logf(msg string, a ...interface{}) {
	fmt.Printf(msg, a...)
	fmt.Println()
}

func main() {
	r := runner.NewRunner()
	err := r.CacheScript("src/scripts/test.js")
	if err != nil {
		panic(err)
	}

	consumer.RedisConnection(func(m *redis.Message) {
		var e runner.Event
		err := json.Unmarshal([]byte(m.Payload), &e)
		if err != nil {
			fmt.Printf("err: %+v\n", err)
			return
		}

		r.Run(&e, func(v *v8go.Value, err error) {
			fmt.Printf("v: %+v\n", v)
			fmt.Printf("err: %+v\n", err)
		})
	})

	ok := make(chan bool, 1)

	<-ok

}
