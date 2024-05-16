package pluginsource

import (
	"fmt"
	"log/slog"

	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/internal/itf"
	"github.com/sandrolain/event-runner/src/internal/utils"
)

type PluginEventOutput struct {
	connection *PluginEventConnection
	slog       *slog.Logger
	config     config.Output
	stopped    bool
}

func (s *PluginEventOutput) Ingest(c <-chan itf.RunnerResult) (err error) {
	go func() {
		for !s.stopped {
			res := <-c
			err := s.send(res)
			if err != nil {
				res.Nak()
				s.slog.Error("error sending data", "err", err)
			} else {
				res.Ack()
			}
		}
	}()
	return
}

func (s *PluginEventOutput) send(result itf.RunnerResult) (err error) {
	data, err := result.Data()
	if err != nil {
		err = fmt.Errorf("error getting data: %w", err)
		return
	}

	cfg, err := result.Config()
	if err != nil {
		err = fmt.Errorf("error getting config: %w", err)
		return
	}

	var serData []byte
	if s.config.Marshal == "" {
		var ok bool
		serData, ok = data.([]byte)
		if !ok {
			err = fmt.Errorf("error casting data: %w", err)
			return
		}
	} else {
		serData, err = utils.Marshal(s.config.Marshal, data)
		if err != nil {
			err = fmt.Errorf("error serializing data: %w", err)
			return
		}
	}

	var uuid string
	if u, ok := cfg["uuid"]; ok && u != "" {
		uuid = cfg["uuid"]
	}

	topic := s.config.Topic
	if t, ok := cfg["topic"]; ok && t != "" {
		topic = t
	}

	meta, err := result.Metadata()
	if err != nil {
		err = fmt.Errorf("error getting metadata: %w", err)
		return
	}

	e := s.connection.plugin.Output(uuid, topic, serData, meta)
	if e != nil {
		err = fmt.Errorf("error sending data: %w", e)
		return
	}

	return
}

func (s *PluginEventOutput) Close() (err error) {
	s.stopped = true
	return
}
