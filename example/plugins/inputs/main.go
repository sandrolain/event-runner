package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	plugin "github.com/sandrolain/event-runner/src/plugin/bootstrap"
	"github.com/sandrolain/event-runner/src/plugin/proto"
)

func main() {
	plugin.Start(plugin.StartOptions{
		Service: &server{},
		Callback: func() error {

			time.Sleep(2 * time.Second)

			plugin.SetReady()

			return nil
		},
	})
}

type server struct {
	proto.UnimplementedPluginServiceServer
}

func (s *server) Status(ctx context.Context, in *proto.StatusReq) (*proto.StatusRes, error) {
	return plugin.GetStatusResponse(), nil
}

func (s *server) Shutdown(ctx context.Context, in *proto.ShutdownReq) (*proto.ShutdownRes, error) {
	return plugin.Shutdown(in.Wait), nil
}

func (s *server) Input(req *proto.InputReq, stream proto.PluginService_InputServer) error {
	var topic string
	for _, v := range req.Configs {
		if v.Name == "topic" {
			topic = v.Value
			break
		}
	}

	timer := time.NewTicker(500 * time.Microsecond)

	for {
		select {
		// Exit on stream context done
		case <-stream.Context().Done():
			return nil
		case <-timer.C:

			resUuid := uuid.New().String()

			res := &proto.InputRes{
				Uuid:  resUuid,
				Topic: topic,
				Data:  []byte(time.Now().Format(time.RFC3339Nano)),
				Metadata: []*proto.Metadata{
					{
						Name:  "uuid",
						Value: resUuid,
					},
					{
						Name:  "topic",
						Value: topic,
					},
					{
						Name:  "time",
						Value: time.Now().String(),
					},
				},
			}

			// Send the Hardware stats on the stream
			err := stream.Send(res)
			if err != nil {
				slog.Error("failed to send hardware stats", "error", err)
			}
		}
	}
}
