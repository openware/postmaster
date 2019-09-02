package consumer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/go-yaml/yaml"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/openware/postmaster/internal/config"
	"github.com/openware/postmaster/pkg/amqp"
	"github.com/openware/postmaster/pkg/env"
	"github.com/openware/postmaster/pkg/eventapi"
)

var (
	Logger = log.Logger
)

func amqpURI() string {
	host := env.FetchDefault("RABBITMQ_HOST", "localhost")
	port := env.FetchDefault("RABBITMQ_PORT", "5672")
	username := env.FetchDefault("RABBITMQ_USERNAME", "guest")
	password := env.FetchDefault("RABBITMQ_PASSWORD", "guest")

	return fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
}

func requireEnvs() {
	env.Must(env.Fetch("SMTP_PASSWORD"))
	env.Must(env.Fetch("SENDER_EMAIL"))
}

// Logger sets settings of zerolog logger.
// Supported environment variables:
// - POSTMASTER_ENV can be either "development" or "production".
// - POSTMASTER_LOG_LEVEL can be "debug", "info", "warn", "error", "fatal", "panic". (default: "debug")
func configureLogger() {
	logLevel, ok := os.LookupEnv("POSTMASTER_LOG_LEVEL")
	if ok {
		level, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			log.Fatal().Err(err)
		}

		zerolog.SetGlobalLevel(level)
		return
	}

	environ, ok := os.LookupEnv("POSTMASTER_ENV")
	if strings.EqualFold("production", environ) {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	env.Logger = Logger
}

func Run(path, tag string) {
	configureLogger()
	requireEnvs()

	conf := new(config.Config)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		Logger.Fatal().Err(err).
			Msgf("can not read file %s", path)
	}

	if err := yaml.Unmarshal([]byte(content), &conf); err != nil {
		Logger.Fatal().Err(err).
			Msgf("can not unmarshal configuration %s", path)
	}

	if err := conf.Validate(); err != nil {
		Logger.Fatal().Err(err).
			Msgf("configuration file %s is not valid", path)
	}

	serveMux := amqp.NewServeMux(amqpURI(), tag, conf.Exchanges, conf.Keychain)
	serveMux.Logger = Logger

	for id := range conf.Events {
		eventConf := conf.Events[id]
		serveMux.HandleFunc(eventConf.Key, eventConf.Exchange, func(event eventapi.RawEvent) {
			Logger.Info().Msgf("processing event %s", eventConf.Key)

			usr, err := eventapi.Unmarshal(event)
			if err != nil {
				Logger.Error().
					Err(err).
					RawJSON("event", event["payload"].([]byte)).
					Msg("can not unmarshal event")
				return
			}

			record := new(eventapi.Record)
			dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				TagName:          "json",
				Result:           &record,
				WeaklyTypedInput: true,
			})

			if err != nil {
				Logger.Error().
					Err(err).
					RawJSON("event", event["payload"].([]byte)).
					Msg("can not unmarshal event")
				return
			}

			if err := dec.Decode(usr.Record); err != nil {
				Logger.Error().
					Err(err).
					RawJSON("event", event["payload"].([]byte)).
					Msg("can not unmarshal event")
				return
			}

			// First language in config is default.
			if record.Language == "" {
				record.Language = conf.Languages[0].Code
			}

			Logger.Info().
				Str("uid", record.User.UID).
				Str("email", record.User.Email).
				Msgf("event received")

			// Checks, that language is supported.
			if !conf.ContainsLanguage(record.Language) {
				Logger.Error().
					Str("language", record.Language).
					Msg("language is not supported")
				return
			}

			if strings.TrimSpace(eventConf.Expression) != "" {
				result, err := expr.Eval(eventConf.Expression, event)
				if err != nil {
					Logger.Error().Err(err).Msg("expression evaluation failed")
				}

				match, ok := result.(bool)
				if !ok {
					Logger.Error().Err(err).Msg("expression result is not boolean")
					return
				}

				logger := Logger.Info().
					Str("uid", record.User.UID).
					Str("email", record.User.Email).
					Interface("match", result)

				if !match {
					logger.Msgf("skipped")
					return
				}

				logger.Msgf("matched")
			}

			tpl := eventConf.Template(record.Language)
			content, err := tpl.Content(event)
			if err != nil {
				Logger.Error().Err(err).Msg("template execution failed")
				return
			}

			email := Email{
				FromAddress: env.Must(env.Fetch("SENDER_EMAIL")),
				FromName:    env.FetchDefault("SENDER_NAME", "postmaster"),
				ToAddress:   record.User.Email,
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
				Logger.Error().Err(err).Msg("failed to send email")
			}
		})
	}

	Logger.Info().Msg("waiting for events")
	if err := serveMux.ListenAndServe(); err != nil {
		Logger.Panic().Err(err).Msg("connection failed")
	}
}
