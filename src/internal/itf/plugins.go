package itf

type EventPlugins interface {
	GetPlugin(id string) (EventPlugin, error)
}

type EventPlugin interface {
	Command(key string) (PluginCommand, error)
}

type PluginCommand interface {
	SetData(any) (err error)
	Exec() (PluginResult, error)
}

type PluginResult interface {
	GetCommand() string
	GetUUID() string
	GetData() (any, error)
	IsAsync() bool
}
