package mailconsumer

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
)

import (
	"github.com/mitchellh/mapstructure"
	"github.com/shal/pigeon/pkg/consumer"
	"github.com/shal/pigeon/pkg/eventapi"
	"github.com/shal/pigeon/pkg/utils"
)

const (
	routingKey = "user.email.confirmation.token"
	exchange   = "barong.events.system"
)
func amqpURI() string {
	host := utils.GetEnv("RABBITMQ_HOST", "localhost")
	port := utils.GetEnv("RABBITMQ_PORT", "5672")
	username := utils.GetEnv("RABBITMQ_USERNAME", "guest")
	password := utils.GetEnv("RABBITMQ_PASSWORD", "guest")

	return fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
}

func procRecord(r eventapi.Event) error {
	// Decode map[string]interface{} to AccountRecord.
	acc := AccountCreatedEvent{}

	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "json",
		Result:           &acc,
		WeaklyTypedInput: true,
	})

	if err != nil {
		return err
	}

	if err := dec.Decode(r); err != nil {
		return err
	}

	tpl, err := template.ParseFiles("templates/sign_up.tpl")
	if err != nil {
		return err
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, acc); err != nil {
		return err
	}

	apiKey := utils.MustGetEnv("SENDGRID_API_KEY")

	email := eventapi.Email{
		FromAddress: utils.GetEnv("SENDER_EMAIL", "noreply@pigeon.com"),
		FromName:    utils.GetEnv("SENDER_NAME", "Pigeon"),
		Subject:     "Confirmation Instructions",
		Reader:      bytes.NewReader(buff.Bytes()),
	}

	if _, err := email.Send(apiKey, acc.User.Email); err != nil {
		return err
	}

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

	go func() {
		for delivery := range deliveries {
			jwtReader, err := eventapi.DeliveryAsJWT(delivery)

			if err != nil {
				log.Println(err)
				return
			}

			jwt, err := ioutil.ReadAll(jwtReader)
			if err != nil {
				log.Println(err)
				return
			}

			log.Printf("Token: %s\n", string(jwt))

			claims, err := eventapi.ParseJWT(string(jwt), eventapi.ValidateJWT)
			if err != nil {
				log.Println(err)
				return
			}

			if err := procRecord(claims.Event); err != nil {
				log.Printf("Consuming: %s\n", err.Error())
			}
		}
	}()

	log.Printf(" [*] Waiting for events. To exit press CTRL+C")
	<-forever
}
