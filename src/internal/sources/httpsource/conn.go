package httpsource

import (
	"fmt"
	"log/slog"
	"net"
	"sync"

	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/internal/itf"
	"github.com/valyala/fasthttp"
)

func NewConnection(cfg config.Connection) (res itf.EventConnection, err error) {
	if cfg.Port == 0 {
		cfg.Port = 8080
	}
	if cfg.Hostname == "" {
		cfg.Hostname = "0.0.0.0"
	}

	addr := fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port)

	conn := &HTTPEventConnection{
		inputs:  make([]*HTTPEventInput, 0),
		config:  cfg,
		slog:    slog.Default().With("context", "HTTP"),
		inputMx: sync.RWMutex{},
		address: addr,
	}

	res = conn

	return
}

type HTTPEventConnection struct {
	listener net.Listener
	inputs   []*HTTPEventInput
	slog     *slog.Logger
	config   config.Connection
	inputMx  sync.RWMutex
	address  string
	started  bool
}

func (c *HTTPEventConnection) startServer() (err error) {
	if c.started {
		return
	}

	slog.Info("starting HTTP server", "addr", c.address)

	// TODO: manage TLS?
	listener, e := net.Listen("tcp", c.address)
	if e != nil {
		err = fmt.Errorf("failed to listen: %w", e)
	}
	c.listener = listener

	e = fasthttp.Serve(listener, func(ctx *fasthttp.RequestCtx) {
		method := string(ctx.Method())
		path := string(ctx.Path())
		found := false
		for _, input := range c.inputs {
			if method == input.config.Method && path == input.config.Topic {
				input.ingest(ctx)
				ctx.SetStatusCode(fasthttp.StatusAccepted)
				found = true
			}
		}
		if !found {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		}
	})
	if e != nil {
		err = fmt.Errorf("failed to start HTTP server: %w", e)
	}
	c.started = false

	return
}

func (c *HTTPEventConnection) NewInput(cfg config.Input) (res itf.EventInput, err error) {
	c.inputMx.Lock()
	defer c.inputMx.Unlock()

	if cfg.Method == "" {
		cfg.Method = fasthttp.MethodPut
	}

	in := &HTTPEventInput{
		connection: c,
		config:     cfg,
		slog:       c.slog.With("topic", cfg.Topic),
		requests:   make(chan *fasthttp.RequestCtx, cfg.Buffer),
	}
	c.startServer()
	c.inputs = append(c.inputs, in)
	res = in
	return
}

func (c *HTTPEventConnection) removeInput(in *HTTPEventInput) (err error) {
	c.inputMx.Lock()
	defer c.inputMx.Unlock()
	for i, v := range c.inputs {
		if v == in {
			c.inputs = append(c.inputs[:i], c.inputs[i+1:]...)
			return
		}
	}
	err = fmt.Errorf("input not found")
	return
}

func (c *HTTPEventConnection) NewOutput(cfg config.Output) (res itf.EventOutput, err error) {
	output := &HTTPEventOutput{
		config: cfg,
		slog:   c.slog,
	}
	output.init()
	res = output
	return
}

func (c *HTTPEventConnection) NewCache(cfg config.Cache) (res itf.EventCache, err error) {
	// TODO
	err = fmt.Errorf("not implemented")
	return
}

func (c *HTTPEventConnection) Close() (err error) {
	c.listener.Close()
	return
}
