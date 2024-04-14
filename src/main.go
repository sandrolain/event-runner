package main

import (
	"log/slog"
	"os"

	"github.com/sandrolain/event-runner/src/internal/setup"
)

func main() {
	env, err := setup.LoadEnvConfig()
	if err != nil {
		slog.Error("error loading env config", "err", err)
		os.Exit(1)
	}

	cfg, err := setup.LoadConfig(env.ConfigPath)
	if err != nil {
		slog.Error("error loading config", "err", err)
		os.Exit(1)
	}

	err = setup.SetupLogger(cfg.Logger)
	if err != nil {
		slog.Error("error setting up logger", "err", err)
		os.Exit(1)
	}

	setup.Exec(cfg)

	setup.HoldExit()
}
