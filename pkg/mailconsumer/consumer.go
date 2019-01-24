package mailconsumer

import (
	"fmt"
	"log"
)

import (
	"github.com/openware/postmaster/pkg/utils"
)

const Exchange = "barong.events.system"

func amqpURI() string {
	host := utils.GetEnv("RABBITMQ_HOST", "localhost")
	port := utils.GetEnv("RABBITMQ_PORT", "5672")
	username := utils.GetEnv("RABBITMQ_USERNAME", "guest")
	password := utils.GetEnv("RABBITMQ_PASSWORD", "guest")

	return fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
}

func Run() {
	amqpURI := amqpURI()

	utils.MustGetEnv("JWT_PUBLIC_KEY")
	utils.MustGetEnv("SENDGRID_API_KEY")

	// Handlers.
	ResetPasswordHandler(amqpURI)
	EmailConfirmationHandler(amqpURI)

	forever := make(chan bool)
	log.Printf(" [*] Waiting for events. To exit press CTRL+C")
	<-forever
}
