package config

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
)

func ApplyDefaults(cfg *Config) (err error) {
	err = defaults.Set(cfg)
	return
}

func Validate(cfg *Config) (err error) {
	err = validator.New(validator.WithRequiredStructEnabled()).Struct(cfg)
	if err != nil {
		return
	}

	if len(cfg.Connections) == 0 {
		err = fmt.Errorf("at least one connection is required")
		return
	}

	if len(cfg.Runners) == 0 {
		err = fmt.Errorf("at least one runner is required")
		return
	}

	if len(cfg.Lines) == 0 {
		err = fmt.Errorf("at least one line is required")
		return
	}

	if len(cfg.Inputs) == 0 {
		err = fmt.Errorf("at least one input is required")
		return
	}

	if len(cfg.Outputs) == 0 {
		err = fmt.Errorf("at least one output is required")
		return
	}

	connIds := make(map[string]bool)
	for _, conn := range cfg.Connections {
		if connIds[conn.ID] {
			err = fmt.Errorf("duplicate connection \"%s\"", conn.ID)
			return
		}
		connIds[conn.ID] = true
	}

	runIds := make(map[string]bool)
	for _, runner := range cfg.Runners {
		if runIds[runner.ID] {
			err = fmt.Errorf("duplicate runner \"%s\"", runner.ID)
			return
		}
		runIds[runner.ID] = true
	}

	lineIds := make(map[string]bool)
	for _, line := range cfg.Lines {
		if lineIds[line.ID] {
			err = fmt.Errorf("duplicate line \"%s\"", line.ID)
			return
		}
		lineIds[line.ID] = true
	}

	inputsIds := make(map[string]bool)
	for _, input := range cfg.Inputs {
		if inputsIds[input.ID] {
			err = fmt.Errorf("duplicate input \"%s\"", input.ID)
			return
		}
		inputsIds[input.ID] = true
	}

	outputsIds := make(map[string]bool)
	for _, output := range cfg.Outputs {
		if outputsIds[output.ID] {
			err = fmt.Errorf("duplicate output \"%s\"", output.ID)
			return
		}
		outputsIds[output.ID] = true
	}

	for _, input := range cfg.Inputs {
		connId := input.ConnectionID
		if connIds[connId] != true {
			err = fmt.Errorf("connection \"%s\" not found for input \"%s\"", connId, input.ID)
			return
		}
	}

	for _, output := range cfg.Outputs {
		connId := output.ConnectionID
		if connIds[connId] != true {
			err = fmt.Errorf("connection \"%s\" not found for output \"%s\"", connId, output.ID)
			return
		}
	}

	for _, line := range cfg.Lines {
		if runIds[line.RunnerID] != true {
			err = fmt.Errorf("runner \"%s\" not found for line \"%s\"", line.RunnerID, line.ID)
			return
		}
		if inputsIds[line.InputID] != true {
			err = fmt.Errorf("input \"%s\" not found for line \"%s\"", line.InputID, line.ID)
			return
		}
		if outputsIds[line.OutputID] != true {
			err = fmt.Errorf("output \"%s\" not found for line \"%s\"", line.OutputID, line.ID)
			return
		}
	}

	return
}
