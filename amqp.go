package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartAmqp[T any](c Config, ctx context.Context, in <-chan T) {

	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
			c.Amqp.Username,
			c.Amqp.Password,
			c.Amqp.Url,
			c.Amqp.Port,
			c.Amqp.Vhost,
		))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	for {
		e, ok := <-in
		if !ok {
			return
		}

		body, err := json.Marshal(e)
		failOnError(err, "Failed to marshal the input event")

		err = ch.PublishWithContext(
			ctx,             // context
			c.Amqp.Exchange, // exchange
			"",              // routing (TODO: Use action type)
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		failOnError(err, "Failed to publish event")
	}
}
