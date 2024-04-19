package es5runner

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
)

type FooCache struct {
	data map[string]any
}

func (f *FooCache) Get(key string) (any, error) {
	return f.data[key], nil
}

func (f *FooCache) Set(key string, data any) error {
	f.data[key] = data
	return nil
}

func (f *FooCache) Del(key string) error {
	delete(f.data, key)
	return nil
}

func (f *FooCache) Close() error {
	return nil
}

func NewCacheWrapper() (res CacheWrapper, vm *goja.Runtime) {
	c := &FooCache{
		data: make(map[string]any),
	}
	vm = goja.New()
	res = CacheWrapper{
		vm:    vm,
		cache: c,
	}
	return
}

func TestCacheWrapper_SetGet(t *testing.T) {
	c, _ := NewCacheWrapper()
	key := "foo"
	value := "bar"

	c.Set(key, value)
	assert.Equal(t, value, c.Get(key))

	value = "baz"
	c.Set(key, value)
	assert.Equal(t, value, c.Get(key))
}

func TestCacheWrapper_Del(t *testing.T) {
	c, _ := NewCacheWrapper()
	key := "foo"
	value := "bar"

	c.Set(key, value)
	assert.Equal(t, value, c.Get(key))

	c.Del(key)
	assert.Equal(t, nil, c.Get(key))
}
