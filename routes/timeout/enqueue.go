package timeout

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"time"

	handle "Gateway/handle"

	amqp "github.com/rabbitmq/amqp091-go"
)

type DiscordUser struct {
	UserId          string   `json:"user_id"`
	UserNickname    string   `json:"user_nickname"`
	UserDisplayName string   `json:"user_display_name"`
	UserRoles       []string `json:"user_roles"`
}

// Request provided in /enqueue
type EnqueueRequest struct {
	Discord struct {
		MessageId string        `json:"message_id"`
		ChannelId string        `json:"channel_id"`
		GuildId   string        `json:"guild_id"`
		Author    DiscordUser   `json:"author"`
		Targets   []DiscordUser `json:"targets"`
	} `json:"discord"`
}

// Create an EnqueueRequest from an HTTP body.
// If the request fails to decode to JSON, an error will be returned.
func Create(body io.ReadCloser) (EnqueueRequest, error) {
	var req EnqueueRequest
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&req)
	return req, err
}

// Handler to process Timeout.Enqueue requests and forward them to an AMQP exchange.
//
// @param c: Channel to intake EnqueueRequests to handle
// @param conn: An AMQP connection
// @param e: AMQP exchange name
func StartTimeoutEnqueue(c <-chan EnqueueRequest, ctx context.Context, conn *amqp.Connection, ex string) {
	// TODO: Handle channel errors
	ch, _ := conn.Channel()
	defer ch.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	for {
		e, ok := <-c
		if !ok {
			return
		}

		// TODO: Check security policy here

		// TODO: Handle marshal errors here
		body, _ := json.Marshal(e)

		err := ch.PublishWithContext(
			ctx,       // context
			ex,        // exchange
			"timeout", // routing
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)

		if err != nil {
			handle.FailOnError(err, "Failed to publish event to AMQP")
		} else {
			log.Printf("Published event from %s@%s to the AMQP broker",
				e.Discord.Author.UserDisplayName,
				e.Discord.Author.UserId)
		}
	}
}
