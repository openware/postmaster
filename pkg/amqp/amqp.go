package amqp

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/streadway/amqp"

	"github.com/openware/postmaster/pkg/eventapi"
)

const (
	MaxRetry = 10
	WaiTime  = 30
)

type muxEntry struct {
	h          Handler
	routingKey string
}

type ServeMux struct {
	Logger zerolog.Logger

	exchanges map[string]string
	keychain  map[string]eventapi.Validator

	tag  string
	addr string
	mu   sync.RWMutex
	m    map[string]map[string]muxEntry

	retries uint8
}

func NewServeMux(addr, tag string, exchanges map[string]string, keychain map[string]eventapi.Validator) *ServeMux {
	return &ServeMux{
		addr:      addr,
		tag:       tag,
		exchanges: exchanges,
		keychain:  keychain,
	}
}

func (mux *ServeMux) declareQueue(channel *amqp.Channel, routingKey string, exchange string) (*amqp.Queue, error) {
	queueName := fmt.Sprintf("postmaster.%s.consumer", routingKey)

	queue, err := channel.QueueDeclare(
		queueName,
		true,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	err = channel.QueueBind(
		queue.Name,
		routingKey,
		exchange,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return &queue, nil
}

func (mux *ServeMux) declareExchange(name string, channel *amqp.Channel) error {
	err := channel.ExchangeDeclare(name, "direct", false, false, false, false, nil)
	if err != nil {
		return err
	}

	return nil
}

func (mux *ServeMux) ListenQueue(
	deliveries <-chan amqp.Delivery,
	handler Handler,
	key, exchangeID string,
) {
	for {
		delivery, ok := <-deliveries
		if !ok {
			mux.Logger.Error().Msgf("stopped listening %s", key)
			return
		}

		mux.Logger.Debug().
			RawJSON("delivery", delivery.Body).
			Msg("delivery received")

		jwtReader, err := eventapi.DeliveryAsJWT(delivery)
		if err != nil {
			mux.Logger.Error().Err(err).Msg("")
			return
		}

		jwt, err := ioutil.ReadAll(jwtReader)
		if err != nil {
			mux.Logger.Error().Err(err).Msg("")
			return
		}

		mux.Logger.Debug().
			Str("token", string(jwt)).
			Msg("token received")

		validator := mux.keychain[exchangeID]
		claims, err := eventapi.ParseJWT(string(jwt), validator.ValidateJWT)
		if err != nil {
			mux.Logger.Debug().
				Str("token", string(jwt)).
				Msg("validation failed")
			return
		}

		handler.ServeAMQP(claims.Event)
	}
}

func (mux *ServeMux) listen() error {
	conn, err := amqp.Dial(mux.addr)
	if err != nil {
		return err
	}

	// Everything is OK with connection.
	mux.Logger.Info().Msgf("successfully connected to %s", mux.addr)
	mux.retries = 1

	notify := conn.NotifyClose(make(chan *amqp.Error))

	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err.Error())
	}

	// Declare exchanges using one channel.
	for id := range mux.m {
		if err != mux.declareExchange(mux.exchanges[id], channel) {
			return fmt.Errorf("exchange: %s", err.Error())
		}
	}

	// Bind queue to exchange and listen.
	for id, events := range mux.m {
		for key, event := range events {
			channel, err := conn.Channel()
			if err != nil {
				return fmt.Errorf("channel: %s", err.Error())
			}

			queue, err := mux.declareQueue(channel, key, mux.exchanges[id])
			if err != nil {
				return fmt.Errorf("queue: %s", err.Error())
			}

			deliveries, err := channel.Consume(queue.Name, mux.tag, true, false, false, false, nil)
			if err != nil {
				return err
			}

			go mux.ListenQueue(deliveries, event.h, event.routingKey, id)
		}
	}

	return <-notify
}

// ListenAndServe listens messages from RabbitMQ.
// Matches special handler for message.
// Tries to establish connection 10 times, one try per 10 second, then returns error.
func (mux *ServeMux) ListenAndServe() error {
	var err error

	for mux.retries <= MaxRetry {
		if mux.retries != 0 {
			mux.Logger.Error().Msgf("trying #%d", mux.retries)
		}

		err = mux.listen()
		mux.Logger.Error().
			Err(err).
			Msg("failed to listen")

		mux.retries += 1
		mux.Logger.Error().Msgf("waiting for %d seconds", WaiTime)
		time.Sleep(WaiTime * time.Second)
	}

	return err
}

func (mux *ServeMux) Handle(routingKey, exchangeID string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if routingKey == "" {
		mux.Logger.Panic().
			Msgf("pattern %s is not valid", routingKey)
	}
	if handler == nil {
		mux.Logger.Panic().
			Msgf("handler with key %s can not be nil ", routingKey)
	}
	if _, exist := mux.m[routingKey]; exist {
		mux.Logger.Panic().
			Msgf("multiple registrations for %s", routingKey)
	}

	if mux.m == nil {
		mux.m = make(map[string]map[string]muxEntry)
	}

	if mux.m[exchangeID] == nil {
		mux.m[exchangeID] = make(map[string]muxEntry)
	}

	mux.m[exchangeID][routingKey] = muxEntry{h: handler, routingKey: routingKey}
}

func (mux *ServeMux) HandleFunc(routingKey, exchangeID string, handler func(raw eventapi.RawEvent)) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if routingKey == "" {
		mux.Logger.Panic().
			Msgf("pattern %s is not valid", routingKey)
	}
	if handler == nil {
		mux.Logger.Panic().
			Msgf("handler with key %s can not be nil ", routingKey)
	}
	if _, exist := mux.m[routingKey]; exist {
		mux.Logger.Panic().
			Msgf("multiple registrations for %s", routingKey)
	}

	if mux.m == nil {
		mux.m = make(map[string]map[string]muxEntry)
	}

	if mux.m[exchangeID] == nil {
		mux.m[exchangeID] = make(map[string]muxEntry)
	}

	mux.m[exchangeID][routingKey] = muxEntry{h: HandlerFunc(handler), routingKey: routingKey}
}

type Handler interface {
	ServeAMQP(raw eventapi.RawEvent)
}

type HandlerFunc func(raw eventapi.RawEvent)

// ServeHTTP calls f(event).
func (f HandlerFunc) ServeAMQP(raw eventapi.RawEvent) {
	f(raw)
}
