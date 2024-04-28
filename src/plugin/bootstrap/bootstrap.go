package plugin

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/eapache/go-resiliency/retrier"
	"github.com/go-playground/validator/v10"
	"github.com/sandrolain/event-runner/src/internal/utils"
	"github.com/sandrolain/event-runner/src/plugin/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type StartOptions struct {
	Service  proto.PluginServiceServer
	Callback func() error
}

type Config struct {
	AppPort   string `env:"APP_PORT" validate:"required"`
	ID        string `env:"PLUGIN_ID" validate:"required"`
	Port      string `env:"PLUGIN_PORT" validate:"required"`
	ConnRetry int    `env:"CONN_RETRY" envDefault:"3" validate:"omitempty"`
	ConnDelay string `env:"CONN_DELAY" envDefault:"1s" validate:"omitempty"`
	Marshal   string `env:"MARSHAL" envDefault:"msgpack" validate:"required,oneof=json msgpack gob"`
}

var cfg Config
var lis net.Listener
var pluginStatus proto.Status = proto.Status_STATUS_STARTUP
var err error
var grpcConn *grpc.ClientConn
var grpcCLient proto.AppServiceClient

func Start(opts StartOptions) {
	e := runStart(opts)
	if e != nil {
		slog.Error("failed to bootstrap", "error", e)
	}
}

func runStart(opts StartOptions) (err error) {
	e := env.Parse(&cfg)
	if e != nil {
		err = fmt.Errorf("cannot parse config: %w", e)
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	e = validate.Struct(cfg)
	if e != nil {
		err = fmt.Errorf("config validation failed: %w", e)
		return
	}

	slog.Info("starting plugin", "port", cfg.Port, "app_port", cfg.AppPort, "id", cfg.ID)

	appAddr := fmt.Sprintf("localhost:%s", cfg.AppPort)

	cRetry := cfg.ConnRetry
	cDelay, e := time.ParseDuration(cfg.ConnDelay)
	if e != nil {
		err = fmt.Errorf("failed to parse conn delay: %w", e)
		return
	}
	r := retrier.New(retrier.ConstantBackoff(cRetry, cDelay), nil)
	e = r.Run(func() (err error) {
		grpcConn, err = grpc.Dial(appAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return
		}
		grpcCLient = proto.NewAppServiceClient(grpcConn)
		return
	})
	if e != nil {
		err = fmt.Errorf("failed to connect to app: %w", e)
		return
	}

	slog.Info("listening on port", "port", cfg.Port)
	lis, e = net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if e != nil {
		err = fmt.Errorf("failed to listen: %w", err)
		return
	}

	// Create a new gRPC server
	s := grpc.NewServer()
	s.RegisterService(&proto.PluginService_ServiceDesc, opts.Service)

	// Start the gRPC server
	errCh := make(chan error)

	go func() {
		slog.Info("serving")
		e := s.Serve(lis)
		if e != nil {
			slog.Error("failed to serve", "error", e)
			errCh <- fmt.Errorf("failed to serve: %w", err)
		}
	}()

	if opts.Callback != nil {
		go func() {
			slog.Info("executing callback")
			if e := opts.Callback(); e != nil {
				slog.Error("callback failed", "error", e)
				SetError(fmt.Errorf("callback failed: %w", err))
			}
		}()
	}

	err = <-errCh

	return
}

func SetError(e error) {
	err = e
	pluginStatus = proto.Status_STATUS_ERROR
}

func SetReady() bool {
	if pluginStatus == proto.Status_STATUS_STARTUP {
		pluginStatus = proto.Status_STATUS_READY
		return true
	}
	return false
}

func GetStatus() proto.Status {
	return pluginStatus
}

func GetStatusResponse() *proto.StatusRes {
	var errMsg *string
	if err != nil {
		m := err.Error()
		errMsg = &m
	}
	// TODO: implememt self-kill after timeout without status requests
	return &proto.StatusRes{
		Status: pluginStatus,
		Error:  errMsg,
	}
}

func Shutdown(delay *string) *proto.ShutdownRes {
	d := "0s"
	if delay != nil {
		d = *delay
	}
	dl, err := time.ParseDuration(d)
	if err != nil {
		slog.Error("failed to parse duration", "error", err)
	}
	pluginStatus = proto.Status_STATUS_SHUTDOWN
	go func() {
		time.Sleep(dl)
		if lis != nil {
			lis.Close()
		}
		os.Exit(0)
	}()
	return &proto.ShutdownRes{}
}

func SuccessResult(in *proto.CommandReq, data any) (res *proto.CommandRes, err error) {
	mData, err := utils.Marshal(cfg.Marshal, data)
	if err != nil {
		err = fmt.Errorf("failed to marshal data: %w", err)
		return
	}
	res = &proto.CommandRes{
		Uuid:    in.Uuid,
		Command: in.Command,
		Result:  proto.Result_RESULT_SUCCESS,
		Data:    mData,
	}
	return
}

func AsyncResult(in *proto.CommandReq) (res *proto.CommandRes, err error) {
	res = &proto.CommandRes{
		Uuid:    in.Uuid,
		Command: in.Command,
		Result:  proto.Result_RESULT_ASYNC,
		Data:    nil,
	}
	return
}

func NotFoundResult(in *proto.CommandReq) (*proto.CommandRes, error) {
	return &proto.CommandRes{
		Uuid:    in.Uuid,
		Command: in.Command,
		Result:  proto.Result_RESULT_ERROR,
		Data:    nil,
	}, status.Errorf(codes.NotFound, "command not found")
}

func SendResult(uuid string, cmd string, result proto.Result, data any) (err error) {
	mData, err := utils.Marshal(cfg.Marshal, data)
	if err != nil {
		err = fmt.Errorf("failed to marshal data: %w", err)
		return
	}
	_, err = grpcCLient.Result(context.Background(), &proto.ResultReq{
		Uuid:    uuid,
		Command: cmd,
		Result:  result,
		Data:    mData,
	})
	return
}

func UnmarshalData[T any](in *proto.CommandReq) (res T, err error) {
	err = utils.Unmarshal(cfg.Marshal, in.Data, &res)
	return
}
