package main

import (
	"fmt"

	"github.com/sandrolain/event-runner/src/runner"
	"rogchap.com/v8go"
)

func main() {
	run := runner.NewRunner()
	err := run.CacheScript("src/scripts/test.js")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		e := &runner.Event{
			Type: "test.js",
			Data: map[string]int{
				"a": 3 * i,
				"b": 4 + i,
			},
		}

		run.Run(e, func(v *v8go.Value, err error) {
			fmt.Printf("v: %+v\n", v)
			fmt.Printf("err: %+v\n", err)
		})
	}

	ok := make(chan bool, 1)

	<-ok

}
