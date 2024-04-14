package natssource

import (
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/internal/itf"
)

func NewConnection(cfg config.Connection) (res itf.EventConnection, err error) {
	var url string
	if cfg.Hostname == "" && cfg.Port == 0 {
		url = nats.DefaultURL
	} else {
		if cfg.Port == 0 {
			cfg.Port = 4222
		}
		if cfg.Hostname == "" {
			cfg.Hostname = "localhost"
		}
		url = fmt.Sprintf("nats://%s:%d", cfg.Hostname, cfg.Port)
	}

	var opts []nats.Option

	if cfg.Token != "" {
		opts = append(opts, nats.Token(cfg.Token))
	} else if cfg.Username != "" && cfg.Password != "" {
		opts = append(opts, nats.UserInfo(cfg.Username, cfg.Password))
	}

	// Connect to a server
	nc, _ := nats.Connect(url, opts...)

	res = &NatsEventConnection{
		connection: nc,
		config:     cfg,
		slog:       slog.Default().With("context", "NATS", "url", url),
	}

	return
}

type NatsEventConnection struct {
	slog       *slog.Logger
	config     config.Connection
	connection *nats.Conn
}

func (c *NatsEventConnection) NewInput(cfg config.Input) (res itf.EventInput, err error) {
	res = &NatsEventInput{
		connection: c.connection,
		config:     cfg,
		slog:       c.slog.With("subject", cfg.Topic, "stream", cfg.Stream),
	}
	return
}

func (c *NatsEventConnection) NewOutput(cfg config.Output) (res itf.EventOutput, err error) {
	res = &NatsEventOutput{
		connection: c.connection,
		config:     cfg,
		slog:       c.slog,
	}
	return
}

func (c *NatsEventConnection) Close() (err error) {
	c.connection.Close()
	return
}
