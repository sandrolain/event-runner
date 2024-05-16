package plugins

import (
	"context"

	"github.com/sandrolain/event-runner/src/plugin/proto"
)

func (p *EventPlugin) Output(uuid string, topic string, data []byte, metadata map[string][]string) (err error) {
	ctx := context.TODO()

	md := make([]*proto.Metadata, 0)
	for k, v := range metadata {
		for _, vv := range v {
			md = append(md, &proto.Metadata{
				Name:  k,
				Value: vv,
			})
		}
	}

	_, err = p.plugin.client.Output(ctx, &proto.OutputReq{
		Uuid:     uuid,
		Topic:    topic,
		Data:     data,
		Metadata: md,
	})

	return
}
