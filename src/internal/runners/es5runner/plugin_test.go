package es5runner

import (
	"fmt"
	"testing"

	"github.com/dop251/goja"
	"github.com/sandrolain/event-runner/src/internal/itf"
	"github.com/stretchr/testify/assert"
)

type TestEventPlugins struct {
	Plugins map[string]itf.EventPlugin
}

func (p *TestEventPlugins) GetPlugin(id string) (res itf.EventPlugin, err error) {
	res, ok := p.Plugins[id]
	if !ok {
		err = fmt.Errorf("plugin %s not found", id)
	}
	return
}

func NewPlugins(plugins map[string]itf.EventPlugin) *TestEventPlugins {
	return &TestEventPlugins{
		Plugins: plugins,
	}
}

type TestEventPlugin struct {
	Result itf.PluginCommandResult
}

func (p *TestEventPlugin) Command(name string, data any) (itf.PluginCommandResult, error) {
	if p.Result.GetCommand() != name {
		return nil, fmt.Errorf("command %s not found", name)
	}

	return p.Result, nil
}

func (p *TestEventPlugin) Input(buffer int, config map[string]string) (<-chan itf.PluginInput, error) {
	return nil, nil
}

func (p *TestEventPlugin) Output(uuid string, topic string, data []byte, metadata map[string][]string) error {
	return nil
}

type TestResult struct {
	Command string
	UUID    string
	Data    any
	Async   bool
}

func (p *TestResult) GetCommand() string {
	return p.Command
}

func (p *TestResult) GetUUID() string {
	return p.UUID
}

func (p *TestResult) GetData() (any, error) {
	return p.Data, nil
}

func (p *TestResult) IsAsync() bool {
	return p.Async
}

func TestPluginsWrapper(t *testing.T) {
	plugin1 := &TestEventPlugin{
		Result: &TestResult{Command: "foo", UUID: "foo", Data: "foo", Async: false},
	}
	plugin2 := &TestEventPlugin{
		Result: &TestResult{Command: "bar", UUID: "bar", Data: "bar", Async: true},
	}
	ins := NewPlugins(map[string]itf.EventPlugin{
		"foo": plugin1,
		"bar": plugin2,
	})
	vm := goja.New()
	plugins := PluginsWrapper{
		vm:      vm,
		plugins: ins,
	}

	plugin := plugins.Get("foo")
	assert.NotNil(t, plugin)
	assert.Equal(t, plugin1, plugin.plugin)

	plugin = plugins.Get("bar")
	assert.NotNil(t, plugin)
	assert.Equal(t, plugin2, plugin.plugin)

	assert.Panics(t, func() {
		plugin = plugins.Get("baz")
	})

	res, err := plugin1.Command("foo", nil)
	assert.Nil(t, err)
	assert.Equal(t, plugin1.Result, res)

	_, err = plugin1.Command("baz", nil)
	assert.NotNil(t, err)

	res, err = plugin2.Command("bar", nil)
	assert.Equal(t, plugin2.Result, res)

	_, err = plugin2.Command("baz", nil)
	assert.NotNil(t, err)
}
