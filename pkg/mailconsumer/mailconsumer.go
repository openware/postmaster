package mailconsumer

import (
	"fmt"
	"github.com/shal/mail-consumer/pkg/utils"
)

func amqpURI() string {
	host := utils.GetEnv("RABBITMQ_HOST", "localhost")
	port := utils.GetEnv("RABBITMQ_PORT", "5672")
	username := utils.GetEnv("RABBITMQ_USERNAME", "guest")
	password := utils.GetEnv("RABBITMQ_PASSWORD", "guest")

	return fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
}

func Run() {
	uri := amqpURI()
	consumer := NewConsumer(uri)
	queue := consumer.DeclareQueue()
	consumer.BindQueue(queue)

	// Listen for events.
	consumer.Run(queue)
}
