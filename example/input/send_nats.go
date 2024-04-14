package main

import (
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	natsClient, _ := nats.Connect(nats.DefaultURL, nats.Token("nats-secret"))
	for {
		time.Sleep(100 * time.Millisecond)
		natsClient.Publish("test.hello", []byte{1, 2, 3, 4})
		slog.Info("sent")
	}
}
