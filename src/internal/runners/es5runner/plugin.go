package es5runner

import (
	"github.com/dop251/goja"
	"github.com/sandrolain/event-runner/src/internal/itf"
)

type PluginsWrapper struct {
	plugins itf.EventPlugins
	vm      *goja.Runtime
}

func (p *PluginsWrapper) Get(key string) (res *PluginWrapper) {
	v, e := p.plugins.GetPlugin(key)
	if e != nil {
		panic(p.vm.NewGoError(e))
	}
	res = &PluginWrapper{
		plugin: v,
		vm:     p.vm,
	}
	return
}

type PluginWrapper struct {
	plugin itf.EventPlugin
	vm     *goja.Runtime
}

func (p *PluginWrapper) Command(name string, data any) (res any, e error) {
	c, e := p.plugin.Command(name)
	if e != nil {
		panic(p.vm.NewGoError(e))
	}
	e = c.SetData(data)
	if e != nil {
		panic(p.vm.NewGoError(e))
	}
	execRes, e := c.Exec()
	if e != nil {
		panic(p.vm.NewGoError(e))
	}
	res, e = execRes.GetData()
	if e != nil {
		panic(p.vm.NewGoError(e))
	}

	return
}
