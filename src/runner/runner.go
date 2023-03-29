package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	v8 "rogchap.com/v8go"
)

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type RunnerCache struct {
	Name   string
	Source string
	Data   *v8.CompilerCachedData
}

type Runner struct {
	Timeout   time.Duration
	CacheData map[string]*RunnerCache
}

func NewRunner() *Runner {
	return &Runner{
		Timeout:   time.Second * 3,
		CacheData: map[string]*RunnerCache{},
	}
}

func (r *Runner) CacheScript(path string) (err error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return
	}

	_, fileName := filepath.Split(path)
	source := string(src)

	iso := v8.NewIsolate()
	v8.NewContext(iso)
	script1, err := iso.CompileUnboundScript(source, fileName, v8.CompileOptions{})
	if err != nil {
		return
	}

	cachedData := script1.CreateCodeCache()
	iso.Dispose()

	r.CacheData[fileName] = &RunnerCache{
		Name:   fileName,
		Source: source,
		Data:   cachedData,
	}

	return
}

func (r *Runner) Run(e *Event, fn func(*v8.Value)) error {
	cache, ok := r.CacheData[e.Type]
	if !ok {
		return fmt.Errorf("script not found: %v", e.Type)
	}

	iso := v8.NewIsolate()
	ctx := v8.NewContext(iso)
	opts := v8.CompileOptions{CachedData: cache.Data}
	script, err := iso.CompileUnboundScript(cache.Source, cache.Name, opts)
	if err != nil {
		return fmt.Errorf("cannot compile script: %v", err)
	}

	jsonData, err := json.Marshal(e.Data)
	if err != nil {
		return fmt.Errorf("cannot marshal data: %v", err)
	}

	vals := make(chan *v8.Value, 1)
	errs := make(chan error, 1)

	go func() {
		_, err := ctx.RunScript(fmt.Sprintf(`const data = %s;`, jsonData), "internal.data.js")
		if err != nil {
			errs <- err
			return
		}
		val, err := script.Run(ctx)
		if err != nil {
			errs <- err
			return
		}
		vals <- val
	}()

	select {
	case val := <-vals:
		fn(val)
	case err := <-errs:
		e := err.(*v8.JSError)
		fmt.Printf("javascript stack trace: %+v", e)
		return err
	case <-time.After(r.Timeout):
		vm := ctx.Isolate()     // get the Isolate from the context
		vm.TerminateExecution() // terminate the execution
		iso.Dispose()
		err := <-errs // will get a termination error back from the running script
		e := err.(*v8.JSError)
		fmt.Printf("javascript stack trace: %+v", e)
		return err
	}
	return nil
}

func ParseEvent(payload *[]byte) (*Event, error) {
	var e Event
	err := json.Unmarshal([]byte(*payload), &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
