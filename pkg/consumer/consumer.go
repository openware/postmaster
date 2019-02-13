package consumer

import (
	"fmt"
	"log"

	"github.com/openware/postmaster/pkg/amqp"
	"github.com/openware/postmaster/pkg/utils"
)

const (
	Exchange = "barong.events.system"
	Tag      = "postmaster"
)

func amqpURI() string {
	host := utils.GetEnv("RABBITMQ_HOST", "localhost")
	port := utils.GetEnv("RABBITMQ_PORT", "5672")
	username := utils.GetEnv("RABBITMQ_USERNAME", "guest")
	password := utils.GetEnv("RABBITMQ_PASSWORD", "guest")

	return fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
}

func Run() {
	amqpURI := amqpURI()

	// List of required environment variables.
	utils.MustGetEnv("JWT_PUBLIC_KEY")
	utils.MustGetEnv("SENDER_EMAIL")
	utils.MustGetEnv("SMTP_PASSWORD")

	serveMux := amqp.NewServeMux(amqpURI, Tag, Exchange)
	serveMux.HandleFunc("user.password.reset.token", ResetPasswordHandler)
	serveMux.HandleFunc("user.email.confirmation.token", EmailConfirmationHandler)

	if err := serveMux.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
