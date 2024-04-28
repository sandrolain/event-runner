package manager

import (
	"context"

	"github.com/sandrolain/event-runner/src/plugin/proto"
)

// implement structure for grpc server for AppService
type AppServiceServer struct {
	proto.UnimplementedAppServiceServer
}

// Result
func (s *AppServiceServer) Result(ctx context.Context, in *proto.ResultReq) (res *proto.ResultRes, err error) {
	return
}

// Input
func (s *AppServiceServer) Input(ctx context.Context, in *proto.InputReq) (res *proto.InputRes, err error) {
	return
}
