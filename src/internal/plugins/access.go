package plugins

import (
	"fmt"

	"github.com/sandrolain/event-runner/src/internal/itf"
)

type EventPlugins struct {
	plugins map[string]*Plugin
}

func NewEventPlugins(plugins []*Plugin) itf.EventPlugins {
	res := EventPlugins{
		plugins: make(map[string]*Plugin),
	}

	for _, plugin := range plugins {
		res.plugins[plugin.Name] = plugin
	}

	return &res
}

func (p *EventPlugins) GetPlugin(id string) (res itf.EventPlugin, err error) {
	plugin, ok := p.plugins[id]
	if !ok {
		fmt.Printf("p.plugins: %+v\n", p.plugins)
		err = fmt.Errorf("plugin %s not found", id)
		return
	}
	res = &EventPlugin{plugin}
	return
}

type EventPlugin struct {
	plugin *Plugin
}
