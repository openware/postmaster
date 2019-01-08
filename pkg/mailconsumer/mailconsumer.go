package mailconsumer

import (
	"fmt"
	"github.com/shal/mail-consumer/pkg/utils"
)

func amqpURI() string {
	//EVENT_API_RABBITMQ_HOST:     localhost
	host := utils.GetEnv("EVENT_API_RABBITMQ_URL", "localhost")
	//EVENT_API_RABBITMQ_PORT:     "5672"
	port := utils.GetEnv("EVENT_API_RABBITMQ_PORT", "5672")
	//EVENT_API_RABBITMQ_USERNAME: guest
	username := utils.GetEnv("EVENT_API_RABBITMQ_USERNAME", "guest")
	//EVENT_API_RABBITMQ_PASSWORD: guest
	password := utils.GetEnv("EVENT_API_RABBITMQ_PASSWORD", "guest")

	return fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
}

func Run() {
	NewConsumer(amqpURI())
}
