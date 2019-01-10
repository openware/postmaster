package mailconsumer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
)

import (
	"github.com/streadway/amqp"
)

const (
	ConsumerTag  = "pigeon"
	RoutingKey   = "account.created"
	ExchangeName = "barong.events.model"
	QueueName    = "pigeon.events.consumer"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
}

func (c *Consumer) BindQueue(queue amqp.Queue) {
	err := c.channel.QueueBind(
		queue.Name,
		RoutingKey,
		ExchangeName,
		false,
		nil,
	)

	if err != nil {
		log.Panicf("Queue Bind: %s", err.Error())
	}
}

func (c *Consumer) DeclareQueue() amqp.Queue {
	err := c.channel.ExchangeDeclare(
		ExchangeName,
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

func parseDelivery(delivery amqp.Delivery) error {
	eventMsg := EventMsg{}

	log.Printf("Delivery: %s\n", delivery.Body)

	// Parse received []byte into JSON.
	if err := json.Unmarshal(delivery.Body, &eventMsg); err != nil {
		log.Printf("Event API Message %s", err.Error())
	}

	// Verify JWT Payload.
	if len(eventMsg.Signatures) < 1 {
		log.Println("")
		return errors.New("no signatures to verify")
	} else if len(eventMsg.Signatures) > 1 {
		return errors.New("multi signature JWT keys does not supported")
	}

	// Build token from received header, payload, signatures.
	tokenStr := fmt.Sprintf("%s.%s.%s",
		eventMsg.Signatures[0].Protected,
		eventMsg.Payload,
		eventMsg.Signatures[0].Signature,
	)

	token, err := jwt.ParseWithClaims(tokenStr, &EventAPIClaims{}, ValidateJWT)

	fmt.Println("Token:", token)

	if err != nil {
		return err
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

	return nil
}

func (c *Consumer) Run(queue amqp.Queue) {
	deliveries, err := c.channel.Consume(
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
		for delivery := range deliveries {
			if err := parseDelivery(delivery); err != nil {
				log.Println(err)
			}
		}
	}()

	log.Printf(" [*] Waiting for events. To exit press CTRL+C")
	<-forever
}

func NewConsumer(amqpURI string) *Consumer {
	// Create a connection.
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
		tag:     ConsumerTag,
	}

	return consumer
}
