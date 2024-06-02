package main

import (
	"context"
	"log"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type DebounceEvent[T any] struct {
	id    string // Unique identifier used for deduplication
	event T      // The actual event to forward or debounced
}

// Provides a debouncer implementation for incoming events
//
// Note: Debouncing requires a dedicated Redis-compliant DB - this is to minimise calls.
// TODO: Explore using alternative debounce methods
//
// @param c: Configuration
// @param ctx: This should be context.Background()
// @param in: Channel to provide events that will be debounced.
// @param out: Debounced events will be emitted here.
func StartDebouncer[T any](c Config, ctx context.Context, in <-chan DebounceEvent[T], out chan<- T) {
	opt, err := redis.ParseURL(c.Redis.Url)
	failOnError(err, "Failed to connect to Redis")

	client := redis.NewClient(opt)
	defer client.Close()

	for {
		e, ok := <-in
		if !ok {
			// Channel is closed, indicates the app is being shut down
			return
		}

		// Debounce only needs to last 30 seconds at most (mutes do not go below 1 minute)
		res, _ := client.SetNX(ctx, e.id, "timeout", time.Duration(time.Second*30)).Result()
		if res {
			out <- e.event
		} else {
			log.Printf("Event {%s} debounced", e.id)
		}
	}
}
