package main

import (
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/lmittmann/tint"
	"github.com/sandrolain/event-runner/src/config"
	es5Runner "github.com/sandrolain/event-runner/src/internal/runners/es5"
	natsSource "github.com/sandrolain/event-runner/src/internal/sources/nats"
)

func main() {
	logLevel := "DEBUG"
	logFormat := "TEXT"

	slogLevel := new(slog.LevelVar)

	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		slogLevel.Set(slog.LevelDebug)
	case "INFO":
		slogLevel.Set(slog.LevelInfo)
	case "WARN":
		slogLevel.Set(slog.LevelWarn)
	case "ERROR":
		slogLevel.Set(slog.LevelError)
	default:
		slogLevel.Set(slog.LevelInfo)
	}

	var handler slog.Handler
	if strings.ToUpper(logFormat) == "JSON" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel, AddSource: true})
	} else {
		handler = tint.NewHandler(os.Stdout, &tint.Options{Level: slogLevel, AddSource: true})
	}
	slog.SetDefault(slog.New(handler))

	conn, err := natsSource.NewConnection(config.Connection{
		Token: "nats-secret",
	})
	if err != nil {
		panic(err)
	}

	source, err := conn.NewInput(config.Input{
		Name: "test.hello",
	})

	runnerMan, err := es5Runner.New(config.Runner{
		ID:          "es5",
		ProgramPath: "./.trash/prog.js",
	})
	if err != nil {
		panic(err)
	}

	c, err := source.Receive(10)
	if err != nil {
		panic(err)
	}

	runner, err := runnerMan.New()
	if err != nil {
		panic(err)
	}

	resC, err := runner.Ingest(c, 10)
	if err != nil {
		panic(err)
	}

	out, err := conn.NewOutput(config.Output{
		Name: "test.response",
	})
	if err != nil {
		panic(err)
	}

	err = out.Ingest(resC)
	if err != nil {
		panic(err)
	}

	// var runner runners.RunnerManager

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh,
		syscall.SIGTERM, // terminate: stopped by `kill -9 PID`
		syscall.SIGINT,  // interrupt: stopped by Ctrl + C
		syscall.SIGHUP,
		syscall.SIGQUIT,
		os.Interrupt,
	)

	<-exitCh
	os.Exit(0)
}

func LogLevel(level string) {
}
