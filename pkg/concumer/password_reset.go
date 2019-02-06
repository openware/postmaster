package consumer

import (
	"bytes"
	"html/template"
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/openware/postmaster/pkg/eventapi"
	"github.com/openware/postmaster/pkg/utils"
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

	tpl, err := template.ParseFiles("templates/reset_password.tpl")
	if err != nil {
		log.Println(err)
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, acc); err != nil {
		log.Println(err)
	}

	email := Email{
		FromAddress: utils.GetEnv("SENDER_EMAIL", "noreply@postmaster.com"),
		FromName:    utils.GetEnv("SENDER_NAME", "postmaster"),
		ToAddress:   acc.User.Email,
		Subject:     "Reset password Instructions",
		Reader:      bytes.NewReader(buff.Bytes()),
	}

	if err := email.Send(); err != nil {
		log.Println(err)
	}
}
