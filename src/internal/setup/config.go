package setup

import (
	"os"

	"github.com/caarlos0/env/v9"
	"github.com/go-playground/validator/v10"
	"github.com/sandrolain/event-runner/src/config"
	"gopkg.in/yaml.v3"
)

type EnvConfig struct {
	ConfigPath string `env:"CONFIG_PATH" envDefault:"./config.yaml"`
}

func LoadEnvConfig() (cfg EnvConfig, err error) {
	err = env.Parse(&cfg)
	if err != nil {
		return
	}

	err = validator.New(validator.WithRequiredStructEnabled()).Struct(cfg)
	if err != nil {
		return
	}

	return
}

func LoadConfig(filePath string) (cfg config.Config, err error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return
	}
	yaml.Unmarshal(yamlFile, &cfg)

	err = config.Validate(&cfg)
	if err != nil {
		return
	}
	return
}
