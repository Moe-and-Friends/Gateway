package main

import (
	"context"
	"log"

	envconfig "github.com/kelseyhightower/envconfig"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type inputEvent struct{}

func main() {

	var c Config
	err := envconfig.Process("Gateway", &c)
	failOnError(err, "Failed to load config")

	var ctx = context.Background()

	in := make(chan DebounceEvent[inputEvent])
	out := make(chan inputEvent)

	go StartDebouncer[inputEvent](c, ctx, in, out)

	// TODO: Add support for incoming HTTP events

}
