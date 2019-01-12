package mailconsumer

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/smtp"
)

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"github.com/shal/pigeon/pkg/consumer"
	"github.com/shal/pigeon/pkg/eventapi"
	"github.com/shal/pigeon/pkg/utils"
	"github.com/streadway/amqp"
)

const (
	RoutingKey = "account.created"
	Exchange   = "barong.events.model"
)

func smtpURI() string {
	host := utils.GetEnv("SMTP_HOST", "localhost")
	port := utils.GetEnv("SMTP_PORT", "25")

	return fmt.Sprintf("%s:%s", host, port)
}

func amqpURI() string {
	host := utils.GetEnv("RABBITMQ_HOST", "localhost")
	port := utils.GetEnv("RABBITMQ_PORT", "5672")
	username := utils.GetEnv("RABBITMQ_USERNAME", "guest")
	password := utils.GetEnv("RABBITMQ_PASSWORD", "guest")

	return fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
}

func DeliveryAsJWT(delivery amqp.Delivery) (string, error) {
	eventMsg := eventapi.Delivery{}

	log.Printf("Delivery: %s\n", delivery.Body)

	// Parse received []byte into JSON.
	if err := json.Unmarshal(delivery.Body, &eventMsg); err != nil {
		log.Printf("Event API Message %s", err.Error())
	}

	// Verify JWT payload.
	if len(eventMsg.Signatures) < 1 {
		return "", errors.New("no signatures to verify")
	} else if len(eventMsg.Signatures) > 1 {
		return "", errors.New("multi signature JWT keys does not supported")
	}

	// Build token from received header, payload, signatures.
	tokenStr := fmt.Sprintf("%s.%s.%s",
		eventMsg.Signatures[0].Protected,
		eventMsg.Payload,
		eventMsg.Signatures[0].Signature,
	)

	log.Printf("Token: %s\n", tokenStr)

	return tokenStr, nil
}

func parseDelivery(delivery amqp.Delivery, callback func(record AccountRecord) error) error {
	tokenStr, err := DeliveryAsJWT(delivery)
	if err != nil {
		return err
	}

	token, err := jwt.ParseWithClaims(tokenStr, &eventapi.Claims{}, eventapi.ValidateJWT)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*eventapi.Claims)
	if !ok || !token.Valid {
		return errors.New("claims: invalid jwt token")
	}

	// Decode map[string]interface{} to AccountRecord.
	acc := AccountRecord{}

	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "json",
		Result:           &acc,
		WeaklyTypedInput: true,
	})

	if err := dec.Decode(claims.Event.Record); err != nil {
		return err
	}

	log.Println(acc)

	// Send email over using standard SMTP package.
	if err := callback(acc); err != nil {
		log.Println(err)
	}

	return nil
}

func Run() {
	// TODO: Check JWT_PUBLIC_KEY to be set on start.
	amqpUri := amqpURI()
	smtpUri := smtpURI()

	c := consumer.New(amqpUri, Exchange, RoutingKey)
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

	callback := func(record AccountRecord) error {
		if cli, err := smtp.Dial(smtpUri); err != nil {
			return err
		} else if err := SendEmail(record, cli); err != nil {
			return err
		}

		return nil
	}

	go func() {
		for delivery := range deliveries {
			if err := parseDelivery(delivery, callback); err != nil {
				log.Println(err)
			}
		}
	}()

	log.Printf(" [*] Waiting for events. To exit press CTRL+C")
	<-forever
}
