package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/sandrolain/event-runner/src/runner"
	"github.com/segmentio/kafka-go"
)

const (
	EVENTS_TOPIC = "incoming"
	RESULT_TOPIC = "result"
)

func main() {
	// to produce messages
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9093", EVENTS_TOPIC, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	//conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	for i := 0; i < 10; i++ {
		data := make(map[string]interface{})
		data["a"] = rand.Intn(100)
		data["b"] = rand.Intn(100)

		evt := runner.Event{
			Type: "test.js",
			Data: data,
		}
		msg, err := json.Marshal(evt)
		if err != nil {
			fmt.Printf("cannot marshal message: %v\n", err)
			continue
		}

		time.Sleep(time.Millisecond * 500)
		_, err = conn.WriteMessages(kafka.Message{Value: msg})
		if err != nil {
			log.Fatal("failed to write messages:", err)
		}
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
