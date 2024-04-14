package setup

import (
	"fmt"

	"github.com/sandrolain/event-runner/src/config"
	"github.com/sandrolain/event-runner/src/internal/itf"
	"github.com/sandrolain/event-runner/src/internal/runners/es5runner"
	"github.com/sandrolain/event-runner/src/internal/sources/httpsource"
	"github.com/sandrolain/event-runner/src/internal/sources/natssource"
)

type Executor struct {
	connections     map[string]itf.EventConnection
	runnersManagers map[string]itf.RunnerManager
}

func Exec(cfg config.Config) (err error) {
	exec := &Executor{
		connections: make(map[string]itf.EventConnection),
	}

	getInput := func(id string) (res config.Input, err error) {
		for _, cfg := range cfg.Inputs {
			if cfg.ID == id {
				return cfg, nil
			}
		}
		err = fmt.Errorf("input \"%s\" not found", id)
		return
	}

	getOutput := func(id string) (res config.Output, err error) {
		for _, cfg := range cfg.Outputs {
			if cfg.ID == id {
				return cfg, nil
			}
		}
		err = fmt.Errorf("output \"%s\" not found", id)
		return
	}

	for _, cfg := range cfg.Connections {
		c, e := NewConnection(cfg)
		if e != nil {
			err = fmt.Errorf("failed to create connection \"%s\": %w", cfg.ID, e)
			return
		}
		exec.connections[cfg.ID] = c
	}

	for _, cfg := range cfg.Runners {
		r, e := NewRunnerManager(cfg)
		if e != nil {
			err = fmt.Errorf("failed to create runner \"%s\": %w", cfg.ID, e)
			return
		}
		exec.runnersManagers[cfg.ID] = r
	}

	for _, cfg := range cfg.Lines {
		conn, ok := exec.connections[cfg.Connection]
		if !ok {
			err = fmt.Errorf("connection \"%s\" not found", cfg.Connection)
			return
		}
		runMan, ok := exec.runnersManagers[cfg.Runner]
		if !ok {
			err = fmt.Errorf("runner \"%s\" not found", cfg.Runner)
			return
		}

		input, e := getInput(cfg.Input)
		if err != nil {
			err = fmt.Errorf("failed to get input: %w", e)
			return
		}

		output, e := getOutput(cfg.Output)
		if err != nil {
			err = fmt.Errorf("failed to get output: %w", e)
			return
		}

		in, e := conn.NewInput(input)
		if err != nil {
			err = fmt.Errorf("failed to create input: %w", e)
			return
		}
		out, e := conn.NewOutput(output)
		if err != nil {
			err = fmt.Errorf("failed to create output: %w", e)
			return
		}

		inC, e := in.Receive()
		if err != nil {
			err = fmt.Errorf("failed to receive input: %w", e)
			return
		}

		run, e := runMan.New()
		if e != nil {
			err = fmt.Errorf("failed to create runner: %w", e)
			return
		}

		ouC, e := run.Ingest(inC)
		if e != nil {
			err = fmt.Errorf("failed to ingest input: %w", e)
			return
		}

		e = out.Ingest(ouC)
		if e != nil {
			err = fmt.Errorf("failed to ingest output: %w", e)
			return
		}
	}

	return
}

func NewConnection(cfg config.Connection) (res itf.EventConnection, err error) {
	switch cfg.Type {
	case "nats":
		return natssource.NewConnection(cfg)
	case "http":
		return httpsource.NewConnection(cfg)
	}
	err = fmt.Errorf("unknown connection type: %s", cfg.Type)
	return
}

func NewRunnerManager(cfg config.Runner) (res itf.RunnerManager, err error) {
	switch cfg.Type {
	case "es5":
		return es5runner.NewRunner(cfg)
	}
	err = fmt.Errorf("unknown runner type: %s", cfg.Type)
	return
}
