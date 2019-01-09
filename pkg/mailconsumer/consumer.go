package mailconsumer

import (
	"github.com/streadway/amqp"
	"log"
)

const (
	consumerTag  = "pigeon"
	routingKey   = "account.created"
	exchangeName = "barong.events.model"
	QueueName    = "pigeon.events.consumer"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	// done    chan error
}

func (c *Consumer) BindQueue(queue amqp.Queue) {
	err := c.channel.QueueBind(
		queue.Name,
		routingKey,
		exchangeName,
		false,
		nil,
	)

	if err != nil {
		log.Panicf("Queue Bind: %s", err.Error())
	}
}

func (c *Consumer) DeclareQueue() amqp.Queue {
	err := c.channel.ExchangeDeclare(
		exchangeName,
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

	queue, err := c.channel.QueueDeclare(
		QueueName,
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

func (c *Consumer) Run(queue amqp.Queue) {
	msgs, err := c.channel.Consume(
		queue.Name,
		c.tag,
		true,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Panicf("Consuming: %s", err.Error())
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("[x] %s", d.Body)

			// Decode JWT Token to JSON.

			// Parse json.

			// Send email.

		}
	}()

	log.Printf(" [*] Waiting for events. To exit press CTRL+C")
	<-forever

}

func NewConsumer(amqpURI string) *Consumer {
	// Make a connection.
	conn, err := amqp.Dial(amqpURI)

	if err != nil {
		log.Panicf("Dial %s", err.Error())
	} else {
		log.Println("Successfully connected to rabbitmq")
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
	}

	return consumer
}
