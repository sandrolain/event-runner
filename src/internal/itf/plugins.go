package itf

import (
	"time"

	"github.com/sandrolain/event-runner/src/plugin/proto"
)

type EventPlugins interface {
	GetPlugin(id string) (EventPlugin, error)
}

type EventPlugin interface {
	Command(key string, data any) (PluginCommandResult, error)
	Input(buffer int, config map[string]string) (<-chan PluginInput, error)
}

type PluginCommandResult interface {
	GetCommand() string
	GetUUID() string
	GetData() (any, error)
	IsAsync() bool
}

type PluginInput interface {
	GetTime() time.Time
	GetInput() *proto.InputRes
}
