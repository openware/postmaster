package eventapi

import (
	"io"
	"io/ioutil"
)

import (
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Email struct {
	FromAddress string
	Subject     string
	Reader      io.Reader
}

func (e Email) Send(apiKey, email string) (*rest.Response, error) {
	from := mail.NewEmail("", e.FromAddress)
	to := mail.NewEmail("", email)

	text, err := ioutil.ReadAll(e.Reader)
	if err != nil {
		return nil, err
	}

	html := mail.NewContent("text/html", string(text))
	message := mail.NewV3MailInit(from, e.Subject, to, html)

	client := sendgrid.NewSendClient(apiKey)
	resp, err := client.Send(message)
	if err != nil {
		return nil, err
	}

	return resp, err
}
