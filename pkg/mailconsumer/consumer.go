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

func procRecord(r eventapi.Record) {
	// Decode map[string]interface{} to AccountRecord.
	acc := AccountRecord{}

	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "json",
		Result:           &acc,
		WeaklyTypedInput: true,
	})

	if err := dec.Decode(r); err != nil {
		log.Println(err)
		return
	}

	tpl, err := template.ParseFiles("templates/sign_up.tpl")
	if err != nil {
		log.Println(err)
		return
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, r); err != nil {
		log.Println(err)
		return
	}

	apiKey := utils.MustGetEnv("SENDGRID_API_KEY")
	email := eventapi.Email{
		FromAddress: utils.GetEnv("SENDER_EMAIL", "example@domain.com"),
		Subject:     "Confirmation Instructions",
		Reader:      bytes.NewReader(buff.Bytes()),
	}

	if _, err := email.Send(acc.Email, apiKey); err != nil {
		log.Println(err)
		return
	}
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
			procRecord(claims.Event.Record)
		}
	}()

	log.Printf(" [*] Waiting for events. To exit press CTRL+C")
	<-forever
}
