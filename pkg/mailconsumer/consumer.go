package mailconsumer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
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
	routingKey = "account.created"
	exchange   = "barong.events.model"
)

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

func ParseJWT(tokenStr string, callback func(record AccountRecord)) error {
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

	// Send email.
	callback(acc)

	return nil
}

func Run() {
	amqpUri := amqpURI()

	utils.MustGetEnv("JWT_PUBLIC_KEY")
	utils.MustGetEnv("SENDGRID_API_KEY")

	c := consumer.New(amqpUri, exchange, routingKey)
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

	callback := func(r AccountRecord) {
		tpl, err := template.ParseFiles("templates/sign_up.tpl")
		if err != nil {
			log.Println(err)
		}

		buff := bytes.Buffer{}
		if err := tpl.Execute(&buff, r); err != nil {
			log.Println(err)
		}

		email := eventapi.Email{
			FromAddress: utils.GetEnv("SENDER_EMAIL", "example@domain.com"),
			Subject:     "Confirmation Instructions",
			Reader:      bytes.NewReader(buff.Bytes()),
		}

		if err := email.Send(r.Email); err != nil {
			log.Println(err)
		}
	}

	go func() {
		for delivery := range deliveries {
			jwtStr, err := DeliveryAsJWT(delivery)

			if err != nil {
				log.Println(err)
			}

			if err := ParseJWT(jwtStr, callback); err != nil {
				log.Println(err)
			}
		}
	}()

	log.Printf(" [*] Waiting for events. To exit press CTRL+C")
	<-forever
}
