package consumer

import (
	"bytes"
	"github.com/mitchellh/mapstructure"
	"github.com/openware/postmaster/pkg/eventapi"
	"github.com/openware/postmaster/pkg/utils"
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

	templatePath := utils.GetEnv("PASSWORD_RESET_TEMPLATE_PATH", "templates/password_reset.tpl")
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println(err)
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, acc); err != nil {
		log.Println(err)
	}

	email := Email{
		FromAddress: utils.MustGetEnv("SENDER_EMAIL"),
		FromName:    utils.GetEnv("SENDER_NAME", "postmaster"),
		ToAddress:   acc.User.Email,
		Subject:     "Reset password Instructions",
		Reader:      bytes.NewReader(buff.Bytes()),
	}

	password := utils.MustGetEnv("SMTP_PASSWORD")
	conf := SMTPConf{
		Host: utils.GetEnv("SMTP_HOST", "smtp.sendgrid.net"),
		Port: utils.GetEnv("SMTP_PORT", "25"),
		Username: utils.GetEnv("SMTP_USER", "apikey"),
		Password: password,
	}

	if err := NewEmailSender(conf, email).Send(); err != nil {
		log.Println(err)
	}
}
