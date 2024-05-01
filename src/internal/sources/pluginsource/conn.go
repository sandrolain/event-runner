package pluginsource

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/internal/itf"
)

func NewConnection(cfg config.Connection, manager itf.EventPlugins) (res itf.EventConnection, err error) {
	plugin, err := manager.GetPlugin(cfg.PluginID)
	if err != nil {
		err = fmt.Errorf("failed to get plugin: %w", err)
		return
	}

	res = &PluginEventConnection{
		plugin:  plugin,
		config:  cfg,
		slog:    slog.Default().With("context", "HTTP"),
		inputMx: sync.RWMutex{},
	}

	return
}

type PluginEventConnection struct {
	plugin  itf.EventPlugin
	slog    *slog.Logger
	config  config.Connection
	inputMx sync.RWMutex
}

func (c *PluginEventConnection) NewInput(cfg config.Input) (res itf.EventInput, err error) {
	c.inputMx.Lock()
	defer c.inputMx.Unlock()

	res = &PluginEventInput{
		connection: c,
		config:     cfg,
		slog:       c.slog.With("topic", cfg.Topic),
	}
	return
}

func (c *PluginEventConnection) NewOutput(cfg config.Output) (res itf.EventOutput, err error) {
	res = &PluginEventOutput{
		config: cfg,
		slog:   c.slog,
	}
	return
}

func (c *PluginEventConnection) NewCache(cfg config.Cache) (res itf.EventCache, err error) {
	// TODO
	err = fmt.Errorf("not implemented")
	return
}

func (c *PluginEventConnection) Close() (err error) {
	// TODO
	err = fmt.Errorf("not implemented")
	return
}
