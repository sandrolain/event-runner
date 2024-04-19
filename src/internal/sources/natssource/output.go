package natssource

import (
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/internal/itf"
)

type NatsEventOutput struct {
	slog    *slog.Logger
	config  config.Output
	nats    *nats.Conn
	stopped bool
}

func (s *NatsEventOutput) Ingest(c chan itf.RunnerResult) (err error) {
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
			// TODO: marshal from config
			serData, err := json.Marshal(data)
			if err != nil {
				res.Nak()
				s.slog.Error("error marshaling data", "err", err)
				continue
			}
			s.slog.Debug("publishing", "subject", s.config.Topic, "size", len(serData))
			msg := nats.NewMsg(s.config.Topic)
			meta, err := res.Metadata()
			if err != nil {
				res.Nak()
				s.slog.Error("error getting metadata", "err", err)
				continue
			}
			for k, v := range meta {
				for _, vv := range v {
					msg.Header.Add(k, vv)
				}
			}
			msg.Data = serData

			err = s.nats.PublishMsg(msg)
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
