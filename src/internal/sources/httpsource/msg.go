package httpsource

import (
	"time"

	"github.com/valyala/fasthttp"
)

type NatsEventMessage struct {
	time    time.Time
	httpCtx *fasthttp.RequestCtx
}

func (m NatsEventMessage) Topic() (string, error) {
	return string(m.httpCtx.Path()), nil
}

func (m NatsEventMessage) ReplyTo() (string, error) {
	return string(m.httpCtx.Request.Header.Referer()), nil
}

func (m NatsEventMessage) Metadata(key string) (res []string, err error) {
	value := string(m.httpCtx.Request.Header.Peek(key))
	return []string{value}, nil
}

func (m NatsEventMessage) Data() ([]byte, error) {
	return m.httpCtx.Request.Body(), nil
}

func (m NatsEventMessage) Time() (time.Time, error) {
	return m.time, nil
}

func (m NatsEventMessage) Ack() error {
	return nil
}

func (m NatsEventMessage) Nak() error {
	return nil
}
