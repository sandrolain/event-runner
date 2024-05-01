package pluginsource

import (
	"time"

	"github.com/sandrolain/event-runner/src/internal/itf"
)

type PluginEventMessage struct {
	input itf.PluginInput
}

func (m PluginEventMessage) Topic() (string, error) {
	return string(m.input.GetInput().Topic), nil
}

func (m PluginEventMessage) ReplyTo() (res string, err error) {
	metas, err := m.Metadata("reply-to")
	if err != nil {
		return
	}
	if len(metas) == 0 {
		return
	}
	res = metas[0]
	return
}

func (m PluginEventMessage) Metadata(key string) (res []string, err error) {
	res = make([]string, 0)
	for _, m := range m.input.GetInput().Metadata {
		if m.Name == key {
			res = append(res, m.Value)
		}
	}
	return
}

func (m PluginEventMessage) Data() ([]byte, error) {
	return m.input.GetInput().Data, nil
}

func (m PluginEventMessage) DataString() (string, error) {
	return string(m.input.GetInput().Data), nil
}

func (m PluginEventMessage) Time() (time.Time, error) {
	return m.input.GetTime(), nil
}

func (m PluginEventMessage) Ack() error {
	return nil
}

func (m PluginEventMessage) Nak() error {
	return nil
}
