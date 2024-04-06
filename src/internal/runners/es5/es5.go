package es5

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/sandrolain/event-runner/src/internal/itf"
)

type Config struct {
	Name    string `validate:"required"`
	Program string `validate:"required"`
}

func New(c Config) (res itf.RunnerManager, err error) {
	prog, err := goja.Compile(c.Name, c.Program, true)
	if err != nil {
		return
	}
	res = &ES5RunnerManager{
		program: prog,
	}
	return
}

type ES5RunnerManager struct {
	program *goja.Program
	runners []itf.Runner
}

func (r *ES5RunnerManager) New() (res itf.Runner, err error) {
	res = &ES5Runner{
		program: r.program,
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
	c       chan itf.EventMessage
	program *goja.Program
	stopped bool
}

func (r *ES5Runner) Ingest(c chan itf.EventMessage, oSize int) (o chan itf.RunnerResult, e chan error, err error) {
	o = make(chan itf.RunnerResult, oSize)
	e = make(chan error)
	go func() {
		for !r.stopped {
			msg := <-c
			res, err := r.run(msg)
			if err != nil {
				msg.Nak()
				e <- err
				continue
			}
			fmt.Printf("res: %v\n", res)
			if res != nil {
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
	err = vm.Set("message", msg)
	if err != nil {
		return
	}
	hasResult := false
	var data any
	metadata := map[string][]string{}
	vm.Set("setMetadata", func(name string, value string) {
		hasResult = true
		fmt.Printf("setMetadata: %+v = %+v\n", name, value)
		metadata[name] = []string{value}
	})
	vm.Set("setData", func(d any) {
		hasResult = true
		fmt.Printf("setData: %+v\n", d)
		data = d
	})
	_, err = vm.RunProgram(r.program)
	if err != nil {
		return
	}
	rpl, err := msg.ReplyTo()
	if err != nil {
		return
	}
	if hasResult {
		res = &ES5RunnerResult{
			message:     msg,
			destination: rpl,
			metadata:    map[string][]string{},
			data:        data,
		}
	}
	return
}

type ES5RunnerResult struct {
	message     itf.EventMessage
	destination string
	metadata    map[string][]string
	data        any
}

func (r ES5RunnerResult) Destination() (string, error) {
	return r.destination, nil
}

func (r ES5RunnerResult) Metadata(key string) (res map[string][]string, err error) {
	res = r.metadata
	return
}

func (r ES5RunnerResult) Data() (any, error) {
	return r.data, nil
}

func (r ES5RunnerResult) Ack() error {
	return r.message.Ack()
}

func (r ES5RunnerResult) Nak() error {
	return r.message.Nak()
}

func (r ES5RunnerResult) Message() itf.EventMessage {
	return r.message
}
