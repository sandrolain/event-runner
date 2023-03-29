package main

import (
	"fmt"
	"io/ioutil"

	"github.com/sandrolain/event-runner/src/runner"
	"github.com/sandrolain/event-runner/src/sources"
	"github.com/sandrolain/event-runner/src/sources/kafka"
	"github.com/sandrolain/event-runner/src/sources/redis"
	"rogchap.com/v8go"
)

func logf(msg string, a ...interface{}) {
	fmt.Printf(msg, a...)
	fmt.Println()
}

const (
	SCRIPTS_PATH = "src/scripts/"
)

func main() {
	files, err := ioutil.ReadDir(SCRIPTS_PATH)
	if err != nil {
		panic(err)
	}

	r := runner.NewRunner()

	for _, file := range files {
		filePath := fmt.Sprintf("%s%s", SCRIPTS_PATH, file.Name())
		err := r.CacheScript(filePath)
		if err != nil {
			panic(err)
		}
	}

	cb := func(cm *sources.ConsumerMessage, err error) error {
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return nil
		}
		e, err := runner.ParseEvent(&cm.Payload)
		if err != nil {
			return fmt.Errorf("cannot parse Payload: %v\n", err)
		}
		return r.Run(e, func(v *v8go.Value) {
			fmt.Printf("v: %+v\n", v)
		})
	}

	ks := kafka.NewSource()
	go ks.NewConsumer(cb)
	rs := redis.NewSource()
	go rs.NewConsumer(cb)

	ok := make(chan bool, 1)
	<-ok
}
