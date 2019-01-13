package mailconsumer

import (
	"bytes"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"log"
)

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/shal/pigeon/pkg/utils"
)

type MailAccountTpl struct {
	Record AccountRecord `json:"record"`
}

func SendEmail(record AccountRecord) error {
	apiKey := utils.MustGetEnv("SENDGRID_API_KEY")

	email := utils.GetEnv("SENDER_EMAIL", "example@domain.com")
	name := utils.GetEnv("SENDER_NAME", "example@domain.com")
	subject := "Confirmation Instructions"
	from := mail.NewEmail(name, email)
	to := mail.NewEmail("", record.Email)

	tpl, err := template.ParseFiles("templates/sign_up.tpl")
	if err != nil {
		return err
	}

	buff := bytes.Buffer{}
	tpl.Execute(&buff, record)

	html := mail.NewContent("text/html", buff.String())
	message := mail.NewV3MailInit(from, subject, to, html)

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
