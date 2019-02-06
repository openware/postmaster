package consumer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"text/template"

	"github.com/openware/postmaster/pkg/utils"
)

type Email struct {
	FromAddress string
	FromName    string
	ToAddress   string
	Subject     string
	Reader      io.Reader
}

// Compatible with any SMTP server, either "Mailcather" or "SendGrid".
func (e Email) Send() error {
	// Password is required.
	password, exist := os.LookupEnv("SMTP_PASSWORD")
	if !exist {
		return errors.New("password is not set")
	}

	host := utils.GetEnv("SMTP_HOST", "smtp.sendgrid.net")
	port := utils.GetEnv("SMTP_PORT", "25")
	username := utils.GetEnv("SMTP_USER", "apikey")
	recipients := []string{e.ToAddress}

	URL := fmt.Sprintf("%s:%s", host, port)

	text, err := ioutil.ReadAll(e.Reader)
	if err != nil {
		return err
	}

	tpl, err := template.ParseFiles("templates/email.tpl")
	if err != nil {
		log.Println(err)
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, e); err != nil {
		log.Println(err)
	}

	msg := append(buff.Bytes(), "\r\n"...)
	msg = append(msg, text...)

	auth := smtp.PlainAuth("", username, password, host)
	if err := smtp.SendMail(URL, auth, e.FromAddress, recipients, msg); err != nil {
		return err
	}

	return nil
}
