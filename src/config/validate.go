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

	connIds := make(map[string]bool)
	for _, conn := range cfg.Connections {
		if connIds[conn.ID] {
			err = fmt.Errorf("connection \"%s\" already exists", conn.ID)
			return
		}
		connIds[conn.ID] = true
	}

	runIds := make(map[string]bool)
	for _, runner := range cfg.Runners {
		if runIds[runner.ID] {
			err = fmt.Errorf("runner \"%s\" already exists", runner.ID)
			return
		}
		runIds[runner.ID] = true
	}

	lineIds := make(map[string]bool)
	for _, line := range cfg.Lines {
		if lineIds[line.ID] {
			err = fmt.Errorf("line \"%s\" already exists", line.ID)
			return
		}
		lineIds[line.ID] = true
	}

	inputsIds := make(map[string]bool)
	for _, input := range cfg.Inputs {
		if inputsIds[input.ID] {
			err = fmt.Errorf("input \"%s\" already exists", input.ID)
			return
		}
		inputsIds[input.ID] = true
	}

	outputsIds := make(map[string]bool)
	for _, output := range cfg.Outputs {
		if outputsIds[output.ID] {
			err = fmt.Errorf("output \"%s\" already exists", output.ID)
			return
		}
		outputsIds[output.ID] = true
	}

	for _, line := range cfg.Lines {
		if connIds[line.Connection] != true {
			err = fmt.Errorf("connection \"%s\" not found for line \"%s\"", line.Connection, line.ID)
			return
		}
		if runIds[line.Runner] != true {
			err = fmt.Errorf("runner \"%s\" not found for line \"%s\"", line.Runner, line.ID)
			return
		}
		if inputsIds[line.Input] != true {
			err = fmt.Errorf("input \"%s\" not found for line \"%s\"", line.Input, line.ID)
			return
		}
		if outputsIds[line.Output] != true {
			err = fmt.Errorf("output \"%s\" not found for line \"%s\"", line.Output, line.ID)
			return
		}
	}

	return
}
