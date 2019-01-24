package consumer

import (
	"log"
)

import (
	"github.com/streadway/amqp"
)

const tag = "postmaster"

type Consumer struct {
	Conn       *amqp.Connection
	Channel    *amqp.Channel
	Tag        string
	RoutingKey string
	Exchange   string
}

func (c *Consumer) BindQueue(queue amqp.Queue) {
	err := c.Channel.QueueBind(
		queue.Name,
		c.RoutingKey,
		c.Exchange,
		false,
		nil,
	)

	if err != nil {
		log.Panicf("Queue Bind: %s", err.Error())
	}
}

func (c *Consumer) DeclareQueue(queueName string) amqp.Queue {
	err := c.Channel.ExchangeDeclare(
		c.Exchange,
		"direct",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Panicf("Exchange: %s", err.Error())
	}

	queue, err := c.Channel.QueueDeclare(
		queueName,
		true,
		true,
		true,
		false,
		nil,
	)

	if err != nil {
		log.Panicf("Queue: %s", err.Error())
	}

	return queue
}

func New(uri, exchange, key string) *Consumer {
	// Create a connection.
	conn, err := amqp.Dial(uri)

	if err != nil {
		log.Panicf("Dial %s", err.Error())
	} else {
		log.Printf("Successfully connected to %s\n", uri)
	}

	// Create a channel.
	channel, err := conn.Channel()

	if err != nil {
		log.Panicf("Channel %s", err.Error())
	}

	consumer := &Consumer{
		Conn:       conn,
		Channel:    channel,
		Tag:        tag,
		Exchange:   exchange,
		RoutingKey: key,
	}

	return consumer
}
