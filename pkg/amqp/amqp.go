package amqp

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/openware/postmaster/pkg/eventapi"
	"github.com/streadway/amqp"
)

const (
	MaxRetry uint8 = 10
	WaiTime = 30
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

	retries uint8
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

func (mux *ServeMux) listen() error {
	conn, err := amqp.Dial(mux.addr)
	if err != nil {
		return err
	}

	// Everything is OK with connection.
	log.Printf("Successfully connected to %s\n", mux.addr)
	mux.retries = 1

	notify := conn.NotifyClose(make(chan *amqp.Error))

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

	log.Printf("Waiting for events...\n")

	connErr := <-notify
	if err := conn.Close(); err != nil {
		log.Println(err)
	}

	return connErr
}
// ListenAndServe listens messages from rabbitmq.
// Matches special handler for message.
// Tries to establish connection 10 times, one try per 10 second, then returns error.
func (mux *ServeMux) ListenAndServe() error {
	var err error

	for mux.retries <= MaxRetry {
		if mux.retries != 0 {
			log.Printf("[ERROR] Try #%d...\n",  mux.retries)
		}

		err = mux.listen()

		log.Printf("[ERROR] %s \n", err)
		mux.retries += 1

		log.Printf("[RECOVER] Sleeping for %d seconds\n", WaiTime)
		time.Sleep(WaiTime * time.Second)
		log.Println("[RECOVER] Awake!")
	}

	return err
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
