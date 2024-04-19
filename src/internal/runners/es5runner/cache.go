package es5runner

import (
	"github.com/dop251/goja"
	"github.com/sandrolain/event-runner/src/internal/itf"
)

type CacheWrapper struct {
	cache itf.EventCache
	vm    *goja.Runtime
}

func (c *CacheWrapper) Get(key string) any {
	v, e := c.cache.Get(key)
	if e != nil {
		panic(c.vm.NewGoError(e))
	}
	return v
}

func (c *CacheWrapper) Set(key string, data any) {
	e := c.cache.Set(key, data)
	if e != nil {
		panic(c.vm.NewGoError(e))
	}
}

func (c *CacheWrapper) Del(key string) {
	e := c.cache.Del(key)
	if e != nil {
		panic(c.vm.NewGoError(e))
	}
}
