package amqp

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/openware/postmaster/pkg/eventapi"
	"github.com/streadway/amqp"
)

type muxEntry struct {
	h          Handler
	routingKey string
}

type ServeMux struct {
	exchange string
	tag      string
	addr     string
	mu       sync.RWMutex
	m        map[string]muxEntry
}

func NewServeMux(addr, tag, exchange string) *ServeMux {
	return &ServeMux{
		addr:     addr,
		tag:      tag,
		exchange: exchange,
	}
}

func (mux *ServeMux) declareQueue(channel *amqp.Channel, routingKey string) (*amqp.Queue, error) {
	queueName := fmt.Sprintf("postmaster.%s.consumer", routingKey)

	queue, err := channel.QueueDeclare(
		queueName,
		true,
		true,
		true,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	err = channel.QueueBind(
		queue.Name,
		routingKey,
		mux.exchange,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return &queue, nil
}

func (mux *ServeMux) declareExchange(channel *amqp.Channel) error {
	err := channel.ExchangeDeclare(
		mux.exchange,
		"direct",
		false,
		false,
		false,
		false,
		nil,
	)

	return err
}

func (mux *ServeMux) declareListener(chann *amqp.Channel, queue amqp.Queue, handler Handler) {
	deliveries, err := chann.Consume(
		queue.Name,
		mux.tag,
		true,
		true,
		false,
		false,
		nil,
	)

	go func() {

		if err != nil {
			log.Panicf("consuming: %s", err.Error())
		}

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

			handler.ServeAMQP(claims.Event)
		}
	}()
}

func (mux *ServeMux) ListenAndServe() error {
	// Create listeners for each mux entry.

	conn, err := amqp.Dial(mux.addr)
	if err != nil {
		log.Panicf("Dial %s", err.Error())
	} else {
		log.Printf("Successfully connected to %s\n", mux.addr)
	}

	// Each event will have own: channel, queue, consumer.
	for k, v := range mux.m {
		channel, err := conn.Channel()

		if err != nil {
			return fmt.Errorf("channel %s", err.Error())
		}

		if err != mux.declareExchange(channel) {
			return fmt.Errorf("exchange %s", err.Error())
		}

		queue, err := mux.declareQueue(channel, k)
		if err != nil {
			return fmt.Errorf("queue: %s", err.Error())
		}

		log.Printf("Listening for %s...\n", k)
		mux.declareListener(channel, *queue, v.h)
	}

	// @Ali: We can recover panics here.

	forever := make(chan bool)
	fmt.Printf("Waiting for events. To exit press CTRL+C")
	<-forever

	return nil
}

func (mux *ServeMux) Handle(routingKey string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if routingKey == "" {
		panic("amqp: invalid pattern")
	}
	if handler == nil {
		panic("amqp: nil handler")
	}
	if _, exist := mux.m[routingKey]; exist {
		panic("amqp: multiple registrations for " + routingKey)
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
	mux.m[routingKey] = muxEntry{h: handler, routingKey: routingKey}
}

func (mux *ServeMux) HandleFunc(routingKey string, handler func(event eventapi.Event)) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if routingKey == "" {
		panic("amqp: invalid pattern")
	}
	if handler == nil {
		panic("amqp: nil handler")
	}
	if _, exist := mux.m[routingKey]; exist {
		panic("amqp: multiple registrations for " + routingKey)
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}

	mux.m[routingKey] = muxEntry{h: HandlerFunc(handler), routingKey: routingKey}
}

type Handler interface {
	ServeAMQP(event eventapi.Event)
}

type HandlerFunc func(eventapi.Event)

// ServeHTTP calls f(event).
func (f HandlerFunc) ServeAMQP(event eventapi.Event) {
	f(event)
}
