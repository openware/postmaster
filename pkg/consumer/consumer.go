package consumer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-yaml/yaml"
	"github.com/openware/postmaster/internal/config"
	"github.com/openware/postmaster/pkg/amqp"
	"github.com/openware/postmaster/pkg/env"
	"github.com/openware/postmaster/pkg/eventapi"
)

func amqpURI() string {
	host := env.FetchDefault("RABBITMQ_HOST", "localhost")
	port := env.FetchDefault("RABBITMQ_PORT", "5672")
	username := env.FetchDefault("RABBITMQ_USERNAME", "guest")
	password := env.FetchDefault("RABBITMQ_PASSWORD", "guest")

	return fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
}

func validate(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Panic(err)
	}

	if _, err := config.Validate(file); err != nil {
		log.Panic(err)
	}
}

func requireEnvs() {
	env.Must(env.Fetch("JWT_PUBLIC_KEY"))
	env.Must(env.Fetch("SMTP_PASSWORD"))
	env.Must(env.Fetch("SENDER_EMAIL"))
}

func Run(path string) {
	conf := config.Config{}

	validate(path)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}

	if err := yaml.Unmarshal([]byte(content), &conf); err != nil {
		log.Panic(err)
	}

	requireEnvs()

	serveMux := amqp.NewServeMux(amqpURI(), conf.AMQP.Tag, conf.AMQP.Exchange)
	for id := range conf.Events {
		eventConf := conf.Events[id]
		serveMux.HandleFunc(eventConf.Key, func(event eventapi.Event) {
			log.Printf("Processing event \"%s\n", eventConf.Key)

			usr, err := eventapi.Unmarshal(event)
			if err != nil {
				log.Println(err)
				return
			}

			// Check, that language is supported.
			if !conf.ContainsLanguage(usr.Language) {
				log.Printf("language %s is not supported", usr.Language)
				return
			}

			tpl := eventConf.Template(usr.Language)
			content, err := tpl.Content(event)
			if err != nil {
				log.Println(err)
				return
			}

			email := Email{
				FromAddress: env.Must(env.Fetch("SENDER_EMAIL")),
				FromName:    env.FetchDefault("SENDER_NAME", "postmaster"),
				ToAddress:   usr.User.Email,
				Subject:     tpl.Subject,
				Reader:      bytes.NewReader(content),
			}

			password := env.Must(env.Fetch("SMTP_PASSWORD"))
			conf := SMTPConf{
				Host:     env.FetchDefault("SMTP_HOST", "smtp.sendgrid.net"),
				Port:     env.FetchDefault("SMTP_PORT", "25"),
				Username: env.FetchDefault("SMTP_USER", "apikey"),
				Password: password,
			}

			if err := NewEmailSender(conf, email).Send(); err != nil {
				log.Println(err)
			}
		})
	}

	if err := serveMux.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
