package plugins

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sandrolain/event-runner/src/internal/itf"
	"github.com/sandrolain/event-runner/src/internal/utils"
	"github.com/sandrolain/event-runner/src/plugin/proto"
)

func (p *EventPlugin) Command(name string, data any) (res itf.PluginCommandResult, err error) {
	d, e := utils.Marshal(p.plugin.Config.Marshal, data)
	if e != nil {
		err = fmt.Errorf("failed to marshal data: %w", e)
		return
	}

	cmdUuid := uuid.New().String()

	p.plugin.slog.Debug("executing command", "name", name, "uuid", cmdUuid)

	cRes, e := p.plugin.client.Command(context.TODO(), &proto.CommandReq{
		Uuid:    cmdUuid,
		Command: name,
		Data:    d,
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

	var resData any
	err = utils.Unmarshal(p.plugin.Config.Marshal, cRes.Data, &resData)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal data: %w", err)
		return
	}

	res = &PluginCommandResult{
		async: async,
		uuid:  cmdUuid,
		data:  resData,
	}
	return
}

type PluginCommandResult struct {
	uuid    string
	command string
	data    any
	async   bool
}

func (r *PluginCommandResult) GetCommand() string {
	return r.command
}

func (r *PluginCommandResult) GetUUID() string {
	return r.uuid
}

func (r *PluginCommandResult) GetData() (any, error) {
	return r.data, nil
}

func (r *PluginCommandResult) IsAsync() bool {
	return r.async
}
