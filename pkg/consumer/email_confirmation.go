package consumer

import (
	"bytes"
	"html/template"
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/openware/postmaster/pkg/env"
	"github.com/openware/postmaster/pkg/eventapi"
)

func EmailConfirmationHandler(event eventapi.Event) {
	acc := EmailConfirmationEvent{}

	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "json",
		Result:           &acc,
		WeaklyTypedInput: true,
	})

	if err != nil {
		log.Println(err)
	}

	if err := dec.Decode(event); err != nil {
		log.Println(err)
	}

	templatePath := env.FetchDefault("SIGN_UP_TEMPLATE_PATH", "templates/sign_up.tpl")
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println(err)
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, acc); err != nil {
		log.Println(err)
	}

	email := Email{
		FromAddress: env.Must(env.Fetch("SENDER_EMAIL")),
		FromName:    env.FetchDefault("SENDER_NAME", "postmaster"),
		ToAddress:   acc.User.Email,
		Subject:     "Email Confirmation Instructions",
		Reader:      bytes.NewReader(buff.Bytes()),
	}

	password := env.Must(env.Fetch("SMTP_PASSWORD"))
	conf := SMTPConf{
		Host: env.FetchDefault("SMTP_HOST", "smtp.sendgrid.net"),
		Port: env.FetchDefault("SMTP_PORT", "25"),
		Username: env.FetchDefault("SMTP_USER", "apikey"),
		Password: password,
	}

	if err := NewEmailSender(conf, email).Send(); err != nil {
		log.Println(err)
	}
}
