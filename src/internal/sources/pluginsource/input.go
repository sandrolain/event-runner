package pluginsource

import (
	"fmt"
	"log/slog"

	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/internal/itf"
)

type PluginEventInput struct {
	connection *PluginEventConnection
	slog       *slog.Logger
	config     config.Input
	c          chan itf.EventMessage
}

func (s *PluginEventInput) Receive() (res <-chan itf.EventMessage, err error) {
	in, e := s.connection.plugin.Input(s.config.Buffer, map[string]string{
		"topic":  s.config.Topic,
		"method": s.config.Method,
		"stream": s.config.Stream,
		"client": s.config.Client,
	})
	if e != nil {
		err = fmt.Errorf("failed to get input: %w", e)
		return
	}

	c := make(chan itf.EventMessage, s.config.Buffer)

	s.c = c
	go func() {
		for r := range in {
			c <- &PluginEventMessage{
				input: r,
			}
		}
	}()
	res = c
	return
}

func (s *PluginEventInput) Close() (err error) {
	if s.c != nil {
		close(s.c)
	}
	return
}
