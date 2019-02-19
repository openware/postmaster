package consumer

import (
	"bytes"
	"github.com/mitchellh/mapstructure"
	"github.com/openware/postmaster/pkg/eventapi"
	"github.com/openware/postmaster/pkg/env"
	"html/template"
	"log"
)

func ResetPasswordHandler(event eventapi.Event) {
	acc := ResetPasswordEvent{}

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

	templatePath := env.FetchDefault("PASSWORD_RESET_TEMPLATE_PATH", "templates/password_reset.tpl")
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
		Subject:     "Reset password Instructions",
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
