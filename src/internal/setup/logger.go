package setup

import (
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
	"github.com/sandrolain/event-runner/src/config"
)

func SetupLogger(cfg config.Logger) (err error) {
	slogLevel, err := parseLogLevel(cfg.Level)
	if err != nil {
		return
	}

	var handler slog.Handler
	if strings.ToUpper(cfg.Format) == "JSON" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel, AddSource: true})
	} else if cfg.Color {
		handler = tint.NewHandler(os.Stdout, &tint.Options{Level: slogLevel, AddSource: true})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel, AddSource: true})
	}

	slog.SetDefault(slog.New(handler))

	slog.Info("logger configured", "level", cfg.Level, "format", cfg.Format, "color", cfg.Color)

	return
}

func parseLogLevel(s string) (slog.Level, error) {
	var level slog.Level
	var err = level.UnmarshalText([]byte(s))
	return level, err
}
