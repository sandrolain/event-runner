package WASM

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"

	"github.com/dop251/goja"
	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/internal/itf"
	"github.com/tetratelabs/wazero"
)

type Config struct {
	Name        string `validate:"required"`
	Program     string `validate:"required"`
	ProgramPath string `validate:"required"`
}

func New(c config.Runner) (res itf.RunnerManager, err error) {
	var program []byte
	if c.ProgramPath != "" {
		program, err = os.ReadFile(c.ProgramPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read program file %s: %w", c.ProgramPath, err)
		}
	} else {
		program, err = base64.StdEncoding.DecodeString(c.ProgramB64)
		if err != nil {
			return nil, fmt.Errorf("unable to decode program: %w", err)
		}
	}
	programContent := string(program)

	prog, err := goja.Compile(c.ID, programContent, true)
	if err != nil {
		return
	}
	res = &WASMRunnerManager{
		program: prog,
	}
	return
}

type WASMRunnerManager struct {
	program *goja.Program
	runners []itf.Runner
}

func (r *WASMRunnerManager) New() (res itf.Runner, err error) {
	res = &WASMRunner{
		slog:    slog.Default().With("context", "WASM"),
		program: r.program,
	}
	r.runners = append(r.runners, res)
	return
}

func (r *WASMRunnerManager) StopAll() error {
	for _, runner := range r.runners {
		err := runner.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}

type WASMRunner struct {
	slog    *slog.Logger
	program *goja.Program
	stopped bool
}

func (r *WASMRunner) Ingest(c chan itf.EventMessage, oSize int) (o chan itf.RunnerResult, err error) {
	o = make(chan itf.RunnerResult, oSize)
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

func (r *WASMRunner) Stop() error {
	r.stopped = true
	return nil
}

func (r *WASMRunner) run(msg itf.EventMessage) (res itf.RunnerResult, err error) {
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
		metadata[name] = []string{value}
	})
	vm.Set("addMetadata", func(name string, value string) {
		hasResult = true
		if metadata[name] == nil {
			metadata[name] = []string{}
		}
		metadata[name] = append(metadata[name], value)
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
		res = &WASMRunnerResult{
			message:     msg,
			destination: rpl,
			metadata:    metadata,
			data:        data,
		}
	}
	return
}

type WASMRunnerResult struct {
	message     itf.EventMessage
	destination string
	metadata    map[string][]string
	data        any
}

func (r WASMRunnerResult) Destination() (string, error) {
	return r.destination, nil
}

func (r WASMRunnerResult) Metadata() (res map[string][]string, err error) {
	res = r.metadata
	return
}

func (r WASMRunnerResult) Data() (any, error) {
	return r.data, nil
}

func (r WASMRunnerResult) Ack() error {
	return r.message.Ack()
}

func (r WASMRunnerResult) Nak() error {
	return r.message.Nak()
}

func (r WASMRunnerResult) Message() itf.EventMessage {
	return r.message
}

func executeWasm(ctx context.Context, wasmData []byte) (err error) {
	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx) // This closes everything this Runtime created.

	// Instantiate WASI, which implements host functions needed for TinyGo to
	// implement `panic`.
	// wasi_snapshot_preview1.MustInstantiate(ctx, r)

	// Instantiate the guest Wasm into the same runtime. It exports the `add`
	// function, implemented in WebAssembly.
	mod, err := r.Instantiate(ctx, wasmData)
	if err != nil {
		err = fmt.Errorf("failed to instantiate: %v", err)
		return
	}

	// alloc := mod.ExportedFunction("alloc")
	// if alloc == nil {
	// 	err = fmt.Errorf("Invalid alloc type")
	// 	return
	// }

	filter := mod.ExportedFunction("filter")
	if filter == nil {
		err = fmt.Errorf("Invalid filter type")
		return
	}

	reqJson, err := gojay.Marshal(req)
	if err != nil {
		return
	}

	reqSize := uint64(len(reqJson))

	allocRes, err := alloc.Call(ctx, reqSize)
	if err != nil {
		err = fmt.Errorf("failed to alloc memory: %v", err)
		return
	}

	reqPtr := allocRes[0]

	if ok := mod.Memory().Write(uint32(reqPtr), reqJson); !ok {
		err = fmt.Errorf("failed to write memory")
		return
	}

	resPtrSize, err := filter.Call(ctx, reqPtr, reqSize)
	if err != nil {
		err = fmt.Errorf("cannot call filter: %v", err)
		return
	}

	resPrt := uint32(resPtrSize[0] >> 32)
	resSize := uint32(resPtrSize[0])

	resJson, ok := mod.Memory().Read(resPrt, resSize)
	if !ok {
		err = fmt.Errorf("cannot read memory")
		return
	}

	return
}
