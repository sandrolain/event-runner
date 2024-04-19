package natssource

import (
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
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

	slog.Info("connecting to NATS server", "url", url)

	// Connect to a server
	nc, _ := nats.Connect(url, opts...)

	res = &NatsEventConnection{
		connection: nc,
		config:     cfg,
		slog:       slog.Default().With("context", "NATS", "url", url),
	}

	return
}

func (c *NatsEventConnection) getJetStream() (res *jetstream.JetStream, err error) {
	if c.js != nil {
		res = c.js
		return
	}
	js, e := jetstream.New(c.connection)
	if e != nil {
		err = fmt.Errorf("error initializing jetstream: %w", e)
		return
	}
	c.js = &js
	res = c.js
	return
}

type NatsEventConnection struct {
	slog       *slog.Logger
	config     config.Connection
	connection *nats.Conn
	js         *jetstream.JetStream
}

func (c *NatsEventConnection) NewInput(cfg config.Input) (res itf.EventInput, err error) {
	res = &NatsEventInput{
		nats:   c.connection,
		config: cfg,
		slog:   c.slog.With("subject", cfg.Topic, "stream", cfg.Stream),
	}
	return
}

func (c *NatsEventConnection) NewOutput(cfg config.Output) (res itf.EventOutput, err error) {
	res = &NatsEventOutput{
		nats:   c.connection,
		config: cfg,
		slog:   c.slog,
	}
	return
}

func (c *NatsEventConnection) NewCache(cfg config.Cache) (res itf.EventCache, err error) {
	js, err := c.getJetStream()
	if err != nil {
		return
	}

	cache := &NatsEventCache{
		nats:   c.connection,
		js:     js,
		config: cfg,
		slog:   c.slog,
	}

	err = cache.init()
	if err != nil {
		return
	}

	res = cache
	return
}

func (c *NatsEventConnection) Close() (err error) {
	c.connection.Close()
	return
}
