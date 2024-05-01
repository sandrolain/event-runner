package plugins

import (
	"context"
	"fmt"
	"time"

	"github.com/sandrolain/event-runner/src/internal/itf"
	"github.com/sandrolain/event-runner/src/plugin/proto"
)

func (p *EventPlugin) Input(buffer int, config map[string]string) (res <-chan itf.PluginInput, err error) {
	cfg := []*proto.Config{}
	for k, v := range config {
		cfg = append(cfg, &proto.Config{
			Name:  k,
			Value: v,
		})
	}

	stream, e := p.plugin.client.Input(context.TODO(), &proto.InputReq{
		Configs: cfg,
	})
	if e != nil {
		err = fmt.Errorf("failed to execute input: %w", e)
		return
	}

	resChan := make(chan itf.PluginInput, buffer)

	go func() {
		for !p.plugin.stopped {
			streamRes, e := stream.Recv()
			if e != nil {
				p.plugin.slog.Error("failed to receive input", "error", e)
				continue
			}
			resChan <- &PluginInput{
				time: time.Now(),
				res:  streamRes,
			}
		}
	}()

	return resChan, nil
}

type PluginInput struct {
	time time.Time
	res  *proto.InputRes
}

func (p *PluginInput) GetTime() time.Time {
	return p.time
}

func (p *PluginInput) GetInput() *proto.InputRes {
	return p.res
}
