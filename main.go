package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	config "Gateway/config"
	handle "Gateway/handle"
	timeout "Gateway/routes/timeout"

	envconfig "github.com/kelseyhightower/envconfig"
	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateAMQP(c config.Config) *amqp.Connection {
	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
			c.Amqp.Username,
			c.Amqp.Password,
			c.Amqp.Url,
			c.Amqp.Port,
			c.Amqp.Vhost,
		))
	handle.FailOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func main() {

	var c config.Config
	err := envconfig.Process("Gateway", &c)
	handle.FailOnError(err, "Failed to load config")

	conn := CreateAMQP(c)
	defer conn.Close()

	var ctx = context.Background()

	// Handle requests to enqueue a timeout event.
	ter := make(chan timeout.EnqueueRequest)
	go timeout.StartTimeoutEnqueue(ter, ctx, conn, c.Amqp.Exchange)
	http.HandleFunc("/timeout/enqueue", func(w http.ResponseWriter, r *http.Request) {
		er, err := timeout.Create(r.Body)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		log.Printf("Processing message sent by %s@%s", er.Discord.Author.UserDisplayName, er.Discord.Author.UserId)
		ter <- er
		w.WriteHeader(201)
	})

	log.Println("Now listening for requests on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
