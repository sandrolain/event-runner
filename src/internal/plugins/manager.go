package plugins

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/plugin/proto"
	"google.golang.org/grpc"
)

func NewPluginManager() (res *PluginManager, err error) {
	port, err := GetFreePort()
	if err != nil {
		err = fmt.Errorf("cannot get free port: %w", err)
		return
	}

	l := slog.Default().With("ctx", "plugin manager")

	l.Info("starting plugin manager", "port", port)

	res = &PluginManager{
		slog:    l,
		port:    port,
		plugins: make(map[string][]*Plugin),
	}
	return
}

type PluginManager struct {
	slog     *slog.Logger
	port     int
	plugins  map[string][]*Plugin
	listener net.Listener
	server   *grpc.Server
}

func (p *PluginManager) Start() (err error) {
	p.slog.Info("starting plugin manager", "port", p.port)

	p.listener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", p.port))
	if err != nil {
		err = fmt.Errorf("failed to listen: %w", err)
		return
	}

	p.slog.Info("starting gRPC server", "port", p.port)

	p.server = grpc.NewServer()
	proto.RegisterAppServiceServer(p.server, &AppServiceServer{})

	go func() {
		e := p.server.Serve(p.listener)
		if e != nil {
			// TODO: handle async error
			err = fmt.Errorf("failed to serve: %w", e)
			return
		}
	}()
	return
}

func (p *PluginManager) Stop() (err error) {
	for _, list := range p.plugins {
		for _, plugin := range list {
			plugin.Stop()
		}
	}
	p.server.Stop()
	return
}

func (p *PluginManager) CreatePlugin(cfg config.Plugin) (res *Plugin, err error) {
	p.slog.Info("creating plugin", "id", cfg.ID)

	host := "localhost"
	port, err := GetFreePort()
	if err != nil {
		err = fmt.Errorf("cannot get free port: %w", err)
		return
	}

	delay, err := time.ParseDuration(cfg.Delay)
	if err != nil {
		err = fmt.Errorf("cannot parse delay: %w", err)
		return
	}

	id := uuid.New().String()

	slog.Info("creating plugin", "id", id, "host", host, "port", port)

	res = &Plugin{
		Config:    cfg,
		AppPort:   p.port,
		ID:        id,
		Host:      host,
		Port:      port,
		Exec:      cfg.Exec,
		Name:      cfg.ID,
		Args:      cfg.Args,
		Env:       cfg.Env,
		ConnRetry: cfg.Retry,
		ConnDelay: delay,
		Output:    cfg.Output,
		slog:      p.slog.With("plugin", cfg.ID, "id", id),
	}

	pluginsList, ok := p.plugins[cfg.ID]
	if !ok {
		pluginsList = make([]*Plugin, 0)
	}
	p.plugins[cfg.ID] = append(pluginsList, res)

	return
}

func (p *PluginManager) GetPlugin(id string) (res *Plugin, err error) {
	pluginsList, ok := p.plugins[id]
	if !ok {
		err = fmt.Errorf("plugin not found")
		return
	}
	if len(pluginsList) == 0 {
		err = fmt.Errorf("plugin not found")
		return
	}
	// TODO: implement weighted
	res = pluginsList[0]
	return
}
