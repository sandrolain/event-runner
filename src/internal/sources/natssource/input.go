package natssource

import (
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/internal/itf"
)

type NatsEventInput struct {
	slog         *slog.Logger
	config       config.Input
	nats         *nats.Conn
	subscription *nats.Subscription
	c            chan itf.EventMessage
}

func (s *NatsEventInput) Receive() (c chan itf.EventMessage, err error) {
	c = make(chan itf.EventMessage, s.config.Buffer)
	s.c = c
	// TODO NATS stream with consumer group
	s.subscription, err = s.nats.Subscribe(s.config.Topic, func(m *nats.Msg) {
		s.slog.Debug("received", "subject", m.Subject, "size", len(m.Data))
		c <- &NatsEventMessage{
			time: time.Now(),
			msg:  m,
		}
	})
	return
}

func (s *NatsEventInput) Close() (err error) {
	if s.c != nil {
		close(s.c)
	}
	err = s.subscription.Unsubscribe()
	if err != nil {
		s.slog.Error("error unsubscribing", "err", err)
	}
	s.nats.Close()
	return
}
