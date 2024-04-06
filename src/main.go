package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sandrolain/event-runner/src/config"
	es5Runner "github.com/sandrolain/event-runner/src/internal/runners/es5"
	natsSource "github.com/sandrolain/event-runner/src/internal/sources/nats"
)

func main() {
	conn, err := natsSource.NewConnection(config.Connection{
		Token: "nats-secret",
	})
	if err != nil {
		panic(err)
	}

	source, err := conn.NewInput(config.Input{
		Name: "test.hello",
	})

	runnerMan, err := es5Runner.New(es5Runner.Config{
		Name: "es5",
		Program: `
			setMetadata("foo", "bar");
			setData(message.data());
			5;
		`,
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

	resC, errC, err := runner.Ingest(c, 10)
	if err != nil {
		panic(err)
	}

	go func() {
		for err := range errC {
			fmt.Printf("Error: %v\n", err)
		}
	}()

	go func() {
		for res := range resC {
			d, e := res.Data()
			fmt.Println(res.Message().Time())
			fmt.Printf("Result: %v %v\n", d, e)
		}
	}()

	// var runner runners.RunnerManager

	exitCh := make(chan os.Signal)
	signal.Notify(exitCh,
		syscall.SIGTERM, // terminate: stopped by `kill -9 PID`
		syscall.SIGINT,  // interrupt: stopped by Ctrl + C
		syscall.SIGHUP,
		syscall.SIGQUIT,
		os.Kill,
		os.Interrupt,
	)

	<-exitCh
	os.Exit(0)
}
