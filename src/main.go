package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sandrolain/event-runner/src/config"
	es5Runner "github.com/sandrolain/event-runner/src/internal/runners/es5"
	"github.com/sandrolain/event-runner/src/internal/setup"
	httpSource "github.com/sandrolain/event-runner/src/internal/sources/http"
	natsSource "github.com/sandrolain/event-runner/src/internal/sources/nats"
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

	natsConn, err := natsSource.NewConnection(config.Connection{
		Token: "nats-secret",
	})
	if err != nil {
		panic(err)
	}

	httpConn, err := httpSource.NewConnection(config.Connection{
		Port: 8080,
	})
	if err != nil {
		panic(err)
	}

	source, err := natsConn.NewInput(config.Input{
		Topic: "test.hello",
	})
	if err != nil {
		panic(err)
	}

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

	out, err := httpConn.NewOutput(config.Output{
		Method: "PUT",
		Topic:  "http://127.0.0.1:8989/test/hello",
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
