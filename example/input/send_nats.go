package main

import (
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	natsClient, _ := nats.Connect(nats.DefaultURL, nats.Token("nats-secret"))
	for {
		time.Sleep(500 * time.Microsecond)
		natsClient.Publish("test.hello", []byte(time.Now().Format(time.RFC3339Nano)))
		slog.Info("sent")
	}
}
