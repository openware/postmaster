package consumer

import (
	"bytes"
	"html/template"
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/openware/postmaster/pkg/eventapi"
	"github.com/openware/postmaster/pkg/utils"
)

func ResetPasswordHandler(event eventapi.Event)() {
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

	tpl, err := template.ParseFiles("templates/reset_password.tpl")
	if err != nil {
		log.Println(err)
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, acc); err != nil {
		log.Println(err)
	}

	apiKey := utils.MustGetEnv("SENDGRID_API_KEY")

	email := eventapi.Email{
		FromAddress: utils.GetEnv("SENDER_EMAIL", "noreply@postmaster.com"),
		FromName:    utils.GetEnv("SENDER_NAME", "postmaster"),
		Subject:     "Reset password Instructions",
		Reader:      bytes.NewReader(buff.Bytes()),
	}

	if _, err := email.Send(apiKey, acc.User.Email); err != nil {
		log.Println(err)
	}
}
