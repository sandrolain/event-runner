package es5runner

import "github.com/sandrolain/event-runner/src/internal/itf"

type CacheWrapper struct {
	cache itf.EventCache
}

func (c *CacheWrapper) Get(key string) any {
	v, e := c.cache.Get(key)
	if e != nil {
		panic(e)
	}
	return v
}

func (c *CacheWrapper) Set(key string, data any) {
	e := c.cache.Set(key, data)
	if e != nil {
		panic(e)
	}
}

func (c *CacheWrapper) Del(key string) {
	e := c.cache.Del(key)
	if e != nil {
		panic(e)
	}
}
