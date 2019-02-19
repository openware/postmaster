package consumer

import (
	"fmt"
	"log"

	"github.com/openware/postmaster/pkg/amqp"
	"github.com/openware/postmaster/pkg/env"
)

const (
	Exchange = "barong.events.system"
	Tag      = "postmaster"
)

func amqpURI() string {
	host := env.FetchDefault("RABBITMQ_HOST", "localhost")
	port := env.FetchDefault("RABBITMQ_PORT", "5672")
	username := env.FetchDefault("RABBITMQ_USERNAME", "guest")
	password := env.FetchDefault("RABBITMQ_PASSWORD", "guest")

	return fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
}

func Run() {
	amqpURI := amqpURI()

	// List of required environment variables.
	env.Must(env.Fetch("JWT_PUBLIC_KEY"))
	env.Must(env.Fetch("SENDER_EMAIL"))
	env.Must(env.Fetch("SMTP_PASSWORD"))

	serveMux := amqp.NewServeMux(amqpURI, Tag, Exchange)
	serveMux.HandleFunc("user.password.reset.token", ResetPasswordHandler)
	serveMux.HandleFunc("user.email.confirmation.token", EmailConfirmationHandler)

	if err := serveMux.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
