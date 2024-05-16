package main

import (
	"context"
	"fmt"
	"time"

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

func (s *server) Output(ctx context.Context, in *proto.OutputReq) (*proto.OutputRes, error) {
	fmt.Printf("Output: %+v\n", in)
	return &proto.OutputRes{}, nil
}
