package config

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validConfig = Config{
	Logger: Logger{
		Level:  "info",
		Format: "text",
		Color:  true,
	},
	Connections: []Connection{{
		ID:       "test",
		Type:     "nats",
		Hostname: "localhost",
		Port:     4222,
	}},
	Runners: []Runner{{
		ID:         "test",
		Type:       "es5",
		ProgramB64: "Y29uc29sZS5sb2coImhlbGxvIikK",
	}, {
		ID:          "test2",
		Type:        "es5",
		ProgramPath: "test.js",
	}},
	Inputs: []Input{{
		ID:           "test",
		ConnectionID: "test",
		Topic:        "test",
	}},
	Outputs: []Output{{
		ID:           "test",
		ConnectionID: "test",
		Topic:        "test",
	}},
	Lines: []Line{{
		ID:       "test",
		RunnerID: "test",
		InputID:  "test",
		OutputID: "test",
	}},
}

func CopyValidConfig(t *testing.T) (res Config) {
	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(validConfig)
	assert.NoError(t, err)
	err = gob.NewDecoder(&buf).Decode(&res)
	assert.NoError(t, err)
	return
}

func TestValidate(t *testing.T) {
	// TODO: add more test cases
	type testCase struct {
		name string
		cfg  Config
		err  bool
	}

	invalidInputConfig := CopyValidConfig(t)
	invalidInputConfig.Inputs[0].ID = ""

	invalidOutputConfig := CopyValidConfig(t)
	invalidOutputConfig.Outputs[0].ID = ""

	invalidRunnerConfig := CopyValidConfig(t)
	invalidRunnerConfig.Runners[0].ID = ""

	testCases := []testCase{
		{
			name: "empty config",
			cfg:  Config{},
			err:  true,
		},
		{
			name: "invalid input config",
			cfg:  invalidInputConfig,
			err:  true,
		},
		{
			name: "invalid output config",
			cfg:  invalidOutputConfig,
			err:  true,
		},
		{
			name: "invalid runner config",
			cfg:  invalidRunnerConfig,
			err:  true,
		},
		{
			name: "valid config",
			cfg:  validConfig,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(&tc.cfg)
			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
