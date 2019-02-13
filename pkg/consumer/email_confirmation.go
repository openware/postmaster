package consumer

import (
	"bytes"
	"html/template"

	"github.com/mitchellh/mapstructure"
	"github.com/openware/postmaster/pkg/eventapi"
	"github.com/openware/postmaster/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func EmailConfirmationHandler(event eventapi.Event) {
	log.WithFields(log.Fields{"event": event}).Debugln("Email Confirmation Event Received")
	acc := EmailConfirmationEvent{}

	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "json",
		Result:           &acc,
		WeaklyTypedInput: true,
	})

	if err != nil {
		log.Errorln(err)
	}

	if err := dec.Decode(event); err != nil {
		log.Errorln(err)
	}

	templatePath := utils.GetEnv("SIGN_UP_TEMPLATE_PATH", "templates/sign_up.tpl")
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Errorln(err)
	}

	log.WithFields(log.Fields{"message": acc}).Debugf("Rendering template: %s\n", templatePath)

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, acc); err != nil {
		log.Errorln(err)
	}

	email := Email{
		FromAddress: utils.MustGetEnv("SENDER_EMAIL"),
		FromName:    utils.GetEnv("SENDER_NAME", "postmaster"),
		ToAddress:   acc.User.Email,
		Subject:     "Email Confirmation Instructions",
		Reader:      bytes.NewReader(buff.Bytes()),
	}

	password := utils.MustGetEnv("SMTP_PASSWORD")
	conf := SMTPConf{
		Host:     utils.GetEnv("SMTP_HOST", "smtp.sendgrid.net"),
		Port:     utils.GetEnv("SMTP_PORT", "25"),
		Username: utils.GetEnv("SMTP_USER", "apikey"),
		Password: password,
	}

	log.WithFields(log.Fields{"email": email, "config": conf}).Debugln("Sending email")
	if err := NewEmailSender(conf, email).Send(); err != nil {
		log.Errorln(err)
	}
}
