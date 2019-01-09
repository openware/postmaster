package mailconsumer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/streadway/amqp"
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
			// TODO: Parse received string into JSON.
			msg := EventMsg{}

			log.Printf("[x] %s", d.Body)

			data := bytes.NewReader(d.Body)

			err := json.NewDecoder(data).Decode(&msg)

			if err != nil {
				log.Printf("Event API Message %s\n", err.Error())
			}

			// TODO: Decode and verify JWT Payload.

			if len(msg.Signatures) < 1 {
				log.Println("No signatures to verify. Skipping...")
				continue
			} else if len(msg.Signatures) > 1 {
				log.Println("Multi Signature JWT keys does not supported. Skipping...")
				continue
			}

			tokenStr := fmt.Sprintf("%s.%s.%s",
				msg.Signatures[0].Protected,
				msg.Payload,
				msg.Signatures[0].Signature,
			)

			token, err := jwt.ParseWithClaims(tokenStr, &EventAPIClaims{}, ValidateJWT)

			fmt.Println("Token:", token)

			if err != nil {
				log.Printf("JWT Parse: %s\n", err.Error())
				continue
			}

			fmt.Print(token.Claims)

			claims, ok := token.Claims.(EventAPIClaims)
			if !ok || !token.Valid {
				fmt.Println(err)
			}

			// TODO: Get email, ... from JWT payload.
			email := claims.Event.Record.Email
			fmt.Println("Email", email)

			// TODO: Send email using SMTP package.
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
