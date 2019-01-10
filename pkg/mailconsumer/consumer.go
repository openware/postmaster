package mailconsumer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shal/pigeon/pkg/utils"
	"log"
)

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"github.com/shal/pigeon/pkg/consumer"
	"github.com/shal/pigeon/pkg/eventapi"
	"github.com/streadway/amqp"
)

const (
	RoutingKey = "account.created"
	Exchange   = "barong.events.model"
)

func amqpURI() string {
	host := utils.GetEnv("RABBITMQ_HOST", "localhost")
	port := utils.GetEnv("RABBITMQ_PORT", "5672")
	username := utils.GetEnv("RABBITMQ_USERNAME", "guest")
	password := utils.GetEnv("RABBITMQ_PASSWORD", "guest")

	return fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
}

func parseDelivery(delivery amqp.Delivery) error {
	eventMsg := eventapi.Delivery{}

	log.Printf("Delivery: %s\n", delivery.Body)

	// Parse received []byte into JSON.
	if err := json.Unmarshal(delivery.Body, &eventMsg); err != nil {
		log.Printf("Event API Message %s", err.Error())
	}

	// Verify JWT payload.
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

	log.Printf("Token: %s\n", tokenStr)

	token, err := jwt.ParseWithClaims(tokenStr, &eventapi.Claims{}, eventapi.ValidateJWT)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*eventapi.Claims)
	if !ok || !token.Valid {
		return errors.New("invalid jwt token")
	}

	// Decode map[string]interface{} to AccountRecord.
	acc := AccountRecord{}
	if err := mapstructure.Decode(claims.Event.Record, &acc); err != nil {
		return err
	}

	log.Println(acc)

	// TODO: Send email using SMTP package.
	// TODO: Discuss, do we need authentication in SMTP or not.

	return nil
}

func Run() {
	uri := amqpURI()
	c := consumer.New(uri, Exchange, RoutingKey)
	queue := c.DeclareQueue()
	c.BindQueue(queue)

	deliveries, err := c.Channel.Consume(
		queue.Name,
		c.Tag,
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
