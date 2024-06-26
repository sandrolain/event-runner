package httpsource

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/internal/itf"
	"github.com/sandrolain/event-runner/src/internal/utils"
	"github.com/valyala/fasthttp"
)

type HTTPEventOutput struct {
	slog    *slog.Logger
	config  config.Output
	stopped bool
	client  *fasthttp.Client
}

func (s *HTTPEventOutput) init() {
	// TODO: from configuration
	readTimeout, _ := time.ParseDuration("500ms")
	writeTimeout, _ := time.ParseDuration("500ms")
	s.client = &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}
}

func (s *HTTPEventOutput) Ingest(c <-chan itf.RunnerResult) (err error) {
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

func (s *HTTPEventOutput) send(result itf.RunnerResult) (err error) {
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

	method := strings.ToUpper(s.config.Method)
	if cfg["method"] != "" {
		method = strings.ToUpper(cfg["method"])
	}
	if method == "" {
		method = "POST"
	}

	url := s.config.Topic
	if cfg["topic"] != "" {
		url = cfg["topic"]
	}

	s.slog.Debug("publishing", "method", method, "url", url, "size", len(serData))

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(method)
	req.Header.Set("Content-Type", "application/json")
	req.SetRequestURI(url)
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

	err = s.client.Do(req, res)
	if err != nil {
		err = fmt.Errorf("error sending request: %w", err)
		return
	}

	if res.StatusCode() > 299 {
		err = fmt.Errorf("non-2XX status code: %d", res.StatusCode())
		return
	}

	s.slog.Debug("published", "status", res.StatusCode())

	return
}

func (s *HTTPEventOutput) Close() (err error) {
	s.stopped = true
	return
}
