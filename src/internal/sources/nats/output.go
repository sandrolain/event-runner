package nats

import (
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/internal/itf"
)

type NatsEventOutput struct {
	slog       *slog.Logger
	config     config.Output
	connection *nats.Conn
	stopped    bool
}

func (s *NatsEventOutput) Receive(c chan itf.RunnerResult) (err error) {
	go func() {
		for !s.stopped {
			res := <-c
			data, err := res.Data()
			if err != nil {
				res.Nak()
				s.slog.Error("error getting data", "err", err)
				continue
			}
			// TODO: check data conversion
			b := data.([]byte)
			err = s.connection.Publish(s.config.Name, b)
			if err != nil {
				res.Nak()
				s.slog.Error("error publishing", "err", err)
				continue
			}
			res.Ack()
		}
	}()
	return
}

func (s *NatsEventOutput) Close() (err error) {
	s.stopped = true
	return
}
