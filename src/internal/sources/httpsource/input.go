package httpsource

import (
	"log/slog"
	"time"

	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/internal/itf"
	"github.com/valyala/fasthttp"
)

type HTTPEventInput struct {
	connection *HTTPEventConnection
	slog       *slog.Logger
	config     config.Input
	c          chan itf.EventMessage
	requests   chan *fasthttp.RequestCtx
}

func (s *HTTPEventInput) ingest(ctx *fasthttp.RequestCtx) {
	s.requests <- ctx
}

func (s *HTTPEventInput) Receive() (c chan itf.EventMessage, err error) {
	c = make(chan itf.EventMessage, s.config.Buffer)
	s.c = c
	go func() {
		for r := range s.requests {
			c <- &HttpEventMessage{
				time:    time.Now(),
				httpCtx: r,
			}
		}
	}()
	return
}

func (s *HTTPEventInput) Close() (err error) {
	if s.c != nil {
		close(s.c)
	}
	err = s.connection.removeInput(s)
	return
}
