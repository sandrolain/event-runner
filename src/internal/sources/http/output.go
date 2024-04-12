package nats

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/internal/itf"
	"github.com/valyala/fasthttp"
)

type HTTPEventOutput struct {
	slog    *slog.Logger
	config  config.Output
	stopped bool
}

func (s *HTTPEventOutput) Ingest(c chan itf.RunnerResult) (err error) {
	go func() {
		for !s.stopped {
			res := <-c
			s.slog.Debug("publishing", "subject", s.config.Name)

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

func (s *HTTPEventOutput) send(result itf.RunnerResult) (err error) {
	data, err := result.Data()
	if err != nil {
		err = fmt.Errorf("error getting data: %w", err)
		return
	}

	// TODO: check data conversion
	// TODO: marshal from config
	serData, err := json.Marshal(data)
	if err != nil {
		err = fmt.Errorf("error serializing data: %w", err)
		return
	}

	client := &fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("POST")
	req.SetRequestURI(s.config.Name)
	req.SetBody(serData)

	meta, err := result.Metadata()

	if err != nil {
		err = fmt.Errorf("error getting metadata: %w", err)
		return
	}
	for k, v := range meta {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	err = client.Do(req, res)
	if err != nil {
		err = fmt.Errorf("error sending request: %w", err)
		return
	}

	if res.StatusCode() > 299 {
		err = fmt.Errorf("non-2XX status code: %d", res.StatusCode())
		return
	}
	return
}

func (s *HTTPEventOutput) Close() (err error) {
	s.stopped = true
	return
}
