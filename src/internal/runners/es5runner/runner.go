package es5runner

import (
	"log/slog"
	"time"

	"github.com/dop251/goja"
	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/internal/itf"
	"github.com/sandrolain/event-runner/src/internal/runners"
)

type Config struct {
	Name        string `validate:"required"`
	Program     string `validate:"required"`
	ProgramPath string `validate:"required"`
}

func NewRunner(c config.Runner) (res itf.RunnerManager, err error) {
	program, err := runners.GetProgramContent(c)
	if err != nil {
		return
	}

	timeout, err := time.ParseDuration(c.Timeout)
	if err != nil {
		return
	}

	prog, err := goja.Compile(c.ID, string(program), true)
	if err != nil {
		return
	}

	res = &ES5RunnerManager{
		config:  c,
		program: prog,
		timeout: timeout,
	}
	return
}

type ES5RunnerManager struct {
	program *goja.Program
	runners []itf.Runner
	config  config.Runner
	timeout time.Duration
}

func (r *ES5RunnerManager) New(cache itf.EventCache, plugins itf.EventPlugins) (res itf.Runner, err error) {
	res = &ES5Runner{
		cache:   cache,
		config:  r.config,
		plugins: plugins,
		slog:    slog.Default().With("context", "ES5"),
		program: r.program,
		timeout: r.timeout,
	}
	r.runners = append(r.runners, res)
	return
}

func (r *ES5RunnerManager) StopAll() error {
	for _, runner := range r.runners {
		err := runner.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}

type ES5Runner struct {
	cache   itf.EventCache
	plugins itf.EventPlugins
	config  config.Runner
	timeout time.Duration
	slog    *slog.Logger
	program *goja.Program
	stopped bool
}

func (r *ES5Runner) Ingest(c <-chan itf.EventMessage) (res <-chan itf.RunnerResult, err error) {
	o := make(chan itf.RunnerResult, r.config.Buffer)
	res = o

	for i := 0; i < r.config.Concurrency; i++ {
		go func() {
			for !r.stopped {
				msg := <-c
				res, err := r.run(msg)
				if err != nil {
					msg.Nak()
					r.slog.Error("error running", "err", err)
					continue
				}
				if res != nil {
					r.slog.Debug("got result", "res", res)
					o <- res
				}
				msg.Ack()
			}
		}()
	}

	go func() {
		for !r.stopped {
			msg := <-c
			res, err := r.run(msg)
			if err != nil {
				msg.Nak()
				r.slog.Error("error running", "err", err)
				continue
			}
			if res != nil {
				r.slog.Debug("got result", "res", res)
				o <- res
			}
			msg.Ack()
		}
	}()
	return
}

func (r *ES5Runner) Stop() error {
	r.stopped = true
	return nil
}

func (r *ES5Runner) run(msg itf.EventMessage) (res itf.RunnerResult, err error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	vm.SetMaxCallStackSize(r.config.MaxStack)

	rpl, err := msg.ReplyTo()
	if err != nil {
		return
	}

	result := &ES5RunnerResult{
		message:     msg,
		destination: rpl,
		metadata:    map[string][]string{},
		data:        nil,
		config:      map[string]string{},
		hasResult:   false,
	}

	err = vm.Set("result", result)
	if err != nil {
		return
	}

	err = vm.Set("cache", &CacheWrapper{
		vm:    vm,
		cache: r.cache,
	})
	if err != nil {
		return
	}

	err = vm.Set("plugin", &PluginsWrapper{
		vm:      vm,
		plugins: r.plugins,
	})
	if err != nil {
		return
	}

	err = vm.Set("message", msg)
	if err != nil {
		return
	}

	// Script timeout
	timer := time.AfterFunc(r.timeout, func() {
		vm.Interrupt("Timeout")
	})
	defer timer.Stop()

	v, err := vm.RunProgram(r.program)
	if err != nil {
		return
	}

	hasResult := result.HasResult()

	if !result.HasResult() && v != nil {
		hasResult = true
		result.data = v
	}

	if hasResult {
		res = result
	}

	return
}

type ES5RunnerResult struct {
	message     itf.EventMessage
	destination string
	metadata    map[string][]string
	data        any
	config      map[string]string
	hasResult   bool
}

func (r *ES5RunnerResult) HasResult() bool {
	return r.hasResult
}

func (r *ES5RunnerResult) SetData(data any) {
	r.data = data
	r.hasResult = true
}

func (r *ES5RunnerResult) AddMetadata(name string, value string) {
	if r.metadata[name] == nil {
		r.metadata[name] = []string{}
	}
	r.metadata[name] = append(r.metadata[name], value)
}

func (r *ES5RunnerResult) SetMetadata(name string, value string) {
	r.metadata[name] = []string{value}
}

func (r *ES5RunnerResult) SetConfig(name string, value string) {
	r.config[name] = value
}

func (r *ES5RunnerResult) Destination() (string, error) {
	return r.destination, nil
}

func (r *ES5RunnerResult) Metadata() (res map[string][]string, err error) {
	res = r.metadata
	return
}

func (r *ES5RunnerResult) Config() (res map[string]string, err error) {
	res = r.config
	return
}

func (r *ES5RunnerResult) Data() (any, error) {
	return r.data, nil
}

func (r *ES5RunnerResult) Ack() error {
	return r.message.Ack()
}

func (r *ES5RunnerResult) Nak() error {
	return r.message.Nak()
}

func (r *ES5RunnerResult) Message() itf.EventMessage {
	return r.message
}
