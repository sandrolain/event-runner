package httpsource

import (
	"time"

	"github.com/valyala/fasthttp"
)

type HttpEventMessage struct {
	time    time.Time
	httpCtx *fasthttp.RequestCtx
}

func (m HttpEventMessage) Topic() (string, error) {
	return string(m.httpCtx.Path()), nil
}

func (m HttpEventMessage) ReplyTo() (string, error) {
	return string(m.httpCtx.Request.Header.Referer()), nil
}

func (m HttpEventMessage) Metadata(key string) (res []string, err error) {
	value := string(m.httpCtx.Request.Header.Peek(key))
	return []string{value}, nil
}

func (m HttpEventMessage) Data() ([]byte, error) {
	return m.httpCtx.Request.Body(), nil
}

func (m HttpEventMessage) DataString() (string, error) {
	return string(m.httpCtx.Request.Body()), nil
}

func (m HttpEventMessage) Time() (time.Time, error) {
	return m.time, nil
}

func (m HttpEventMessage) Ack() error {
	return nil
}

func (m HttpEventMessage) Nak() error {
	return nil
}
