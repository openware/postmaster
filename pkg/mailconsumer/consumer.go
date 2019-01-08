package mailconsumer

import (
	"github.com/streadway/amqp"
	"log"
)

const (
	consumerTag = "pigeon"
)

// Exchange Name: barong.events.model

// event: {
//   name: "model.account.created",
//   record: {
//     uid: 'ID092B2AF8E87',
//     email: 'email@example.com',
//     level: 0,
//     otp_enabled: false,
//     confirmation_sent_at: '2018-04-12T17:16:06+03:00',
//     state: 'pending',
//     created_at: '2018-04-12T17:16:06+03:00',
//     updated_at: '2018-04-12T17:16:06+03:00'
//   }
// }

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func (c *Consumer) DeclareQueue() {
	// Declare a queue.
	queue, err := c.channel.QueueDeclare(
		"barong.events.model",
		true,
		true,
		true,
		false,
		nil,
	)

	if err != nil {
		log.Panicf("Channel %s", err.Error())
	}
}

func NewConsumer(amqpURI string) *Consumer {
	// Make a connection.
	conn, err := amqp.Dial(amqpURI)

	if err != nil {
		log.Panicf("Dial %s", err.Error())
	}

	// Create a channel.
	channel, err := conn.Channel()

	if err != nil {
		log.Panicf("Channel %s", err.Error())
	}

	consumer := &Consumer{
		conn:    conn,
		channel: channel,
		tag:     consumerTag,
		done:    make(chan error),
	}

	return consumer
}
