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
	caches          map[string]itf.EventCache
}

func Exec(cfg config.Config) (err error) {
	exec := &Executor{
		connections:     make(map[string]itf.EventConnection),
		runnersManagers: make(map[string]itf.RunnerManager),
		caches:          make(map[string]itf.EventCache),
	}

	getInputConfig := func(id string) (res config.Input, err error) {
		for _, cfg := range cfg.Inputs {
			if cfg.ID == id {
				return cfg, nil
			}
		}
		err = fmt.Errorf("input \"%s\" not found", id)
		return
	}

	getOutputConfig := func(id string) (res config.Output, err error) {
		for _, cfg := range cfg.Outputs {
			if cfg.ID == id {
				return cfg, nil
			}
		}
		err = fmt.Errorf("output \"%s\" not found", id)
		return
	}

	getCache := func(id string) (res itf.EventCache, err error) {
		res, ok := exec.caches[id]
		if ok {
			return
		}

		for _, cfg := range cfg.Caches {
			if cfg.ID == id {
				cacheConn, ok := exec.connections[cfg.ConnectionID]
				if !ok {
					err = fmt.Errorf("connection \"%s\" for cache \"%s\" not found", cfg.ConnectionID, cfg.ID)
					return
				}
				res, err = cacheConn.NewCache(cfg)
				if err != nil {
					err = fmt.Errorf("failed to create cache \"%s\": %w", cfg.ID, err)
					return
				}
				exec.caches[cfg.ID] = res
				return
			}
		}
		err = fmt.Errorf("cache \"%s\" not found", id)
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
		runMan, ok := exec.runnersManagers[cfg.RunnerID]
		if !ok {
			err = fmt.Errorf("runner \"%s\" not found", cfg.RunnerID)
			return
		}

		input, e := getInputConfig(cfg.InputID)
		if err != nil {
			err = fmt.Errorf("failed to get input: %w", e)
			return
		}

		output, e := getOutputConfig(cfg.OutputID)
		if err != nil {
			err = fmt.Errorf("failed to get output: %w", e)
			return
		}

		var cache itf.EventCache
		if cfg.CacheID != "" {
			cache, e = getCache(cfg.CacheID)
			if err != nil {
				err = fmt.Errorf("failed to get cache: %w", e)
				return
			}
		} else {
			// TODO default in-memory cache
		}

		inputConn, ok := exec.connections[input.ConnectionID]
		if !ok {
			err = fmt.Errorf("connection \"%s\" not found", input.ConnectionID)
			return
		}

		outputConn, ok := exec.connections[output.ConnectionID]
		if !ok {
			err = fmt.Errorf("connection \"%s\" not found", output.ConnectionID)
			return
		}

		in, e := inputConn.NewInput(input)
		if err != nil {
			err = fmt.Errorf("failed to create input: %w", e)
			return
		}

		out, e := outputConn.NewOutput(output)
		if err != nil {
			err = fmt.Errorf("failed to create output: %w", e)
			return
		}

		inC, e := in.Receive()
		if err != nil {
			err = fmt.Errorf("failed to receive input: %w", e)
			return
		}

		run, e := runMan.New(cache)
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
