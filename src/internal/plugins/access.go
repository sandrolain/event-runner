package manager

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sandrolain/event-runner/src/internal/itf"
	"github.com/sandrolain/event-runner/src/internal/utils"
	"github.com/sandrolain/event-runner/src/plugin/proto"
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

func (p *EventPlugin) Command(name string) (itf.PluginCommand, error) {
	return &PluginCommand{
		name:   name,
		plugin: p.plugin,
	}, nil
}

type PluginCommand struct {
	plugin *Plugin
	name   string
	data   []byte
}

func (p *PluginCommand) SetData(data any) (err error) {
	d, e := utils.Marshal(p.plugin.Config.Marshal, data)
	if e != nil {
		err = fmt.Errorf("failed to marshal data: %w", e)
		return
	}
	p.data = d
	return
}

func (p *PluginCommand) Exec() (res itf.PluginResult, err error) {
	cmdUuid := uuid.New().String()

	p.plugin.slog.Debug("executing command", "name", p.name, "uuid", cmdUuid)

	cRes, e := p.plugin.client.Command(context.TODO(), &proto.CommandReq{
		Uuid:    cmdUuid,
		Command: p.name,
		Data:    p.data,
	})
	if e != nil {
		err = fmt.Errorf("failed to execute command: %w", e)
		return
	}

	if cRes.Result == proto.Result_RESULT_ERROR {
		err = fmt.Errorf("failed to execute command: %s", cRes.Data)
		return
	}

	async := cRes.Result == proto.Result_RESULT_ASYNC

	if async {
		// TODO: implement
	}

	var data any
	err = utils.Unmarshal(p.plugin.Config.Marshal, cRes.Data, &data)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal data: %w", err)
		return
	}

	fmt.Printf("data: %v\n", data)

	res = &PluginResult{
		async: async,
		uuid:  cmdUuid,
		data:  data,
	}

	return
}

type PluginResult struct {
	uuid    string
	command string
	data    any
	async   bool
}

func (r *PluginResult) GetCommand() string {
	return r.command
}

func (r *PluginResult) GetUUID() string {
	return r.uuid
}

func (r *PluginResult) GetData() (any, error) {
	return r.data, nil
}

func (r *PluginResult) IsAsync() bool {
	return r.async
}
