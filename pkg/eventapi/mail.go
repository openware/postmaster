package eventapi

import (
	"io"
	"io/ioutil"
	"log"
)

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/shal/pigeon/pkg/utils"
)

type Email struct {
	FromAddress string
	Subject     string
	Reader      io.Reader
}

func (e Email) Send(email string) error {
	apiKey := utils.MustGetEnv("SENDGRID_API_KEY")

	from := mail.NewEmail("", e.FromAddress)
	to := mail.NewEmail("", email)

	text, err := ioutil.ReadAll(e.Reader)
	if err != nil {
		return err
	}

	html := mail.NewContent("text/html", string(text))
	message := mail.NewV3MailInit(from, e.Subject, to, html)

	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		return err
	}

	log.Println("Status Code: ", response.StatusCode)
	log.Println("Response Body: ", response.Body)
	log.Println("Response Headers: ", response.Headers)

	return nil
}
