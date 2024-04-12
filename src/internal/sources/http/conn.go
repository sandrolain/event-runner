package nats

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/internal/itf"
	"github.com/valyala/fasthttp"
)

func NewConnection(cfg config.Connection) (res itf.EventConnection, err error) {
	port := cfg.Port
	if port == 0 {
		port = 8080
	}

	conn := &HTTPEventConnection{
		inputs: make([]*HTTPEventInput, 0),
		config: cfg,
		slog:   slog.Default().With("context", "HTTP", "port", port),
	}

	// Create server
	// TODO: manage TLS?
	e := fasthttp.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsPut() {
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			return
		}
		path := string(ctx.Path())
		for _, input := range conn.inputs {
			if path == input.config.Name {
				input.ingest(ctx)
				ctx.SetStatusCode(fasthttp.StatusAccepted)
				return
			}
		}
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	})
	if e != nil {
		err = fmt.Errorf("failed to start server: %w", e)
		return
	}

	res = conn

	return
}

type HTTPEventConnection struct {
	inputs  []*HTTPEventInput
	slog    *slog.Logger
	config  config.Connection
	inputMx sync.RWMutex
}

func (c *HTTPEventConnection) NewInput(cfg config.Input) (res itf.EventInput, err error) {
	c.inputMx.Lock()
	defer c.inputMx.Unlock()
	in := &HTTPEventInput{
		connection: c,
		config:     cfg,
		slog:       c.slog.With("path", cfg.Name),
		requests:   make(chan *fasthttp.RequestCtx, 10),
	}
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
	res = &HTTPEventOutput{
		config: cfg,
		slog:   c.slog.With("url", cfg.Name),
	}
	return
}

func (c *HTTPEventConnection) Close() (err error) {
	//TODO: is possible to end the server?
	return
}
