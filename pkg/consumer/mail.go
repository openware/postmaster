package consumer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/smtp"
	"strings"
	"text/template"
)

type Email struct {
	FromAddress string
	FromName    string
	ToAddress   string
	Subject     string
	Reader      io.Reader
}

type SMTPConf struct {
	Username string
	Password string
	Host     string
	Port     string
}

func (conf SMTPConf) URL() string {
	return fmt.Sprintf("%s:%s", conf.Host, conf.Port)
}

type EmailSender struct {
	conf  *SMTPConf
	email *Email
	send  func(string, smtp.Auth, string, []string, []byte) error
}

func NewEmailSender(conf SMTPConf, email Email) *EmailSender {
	return &EmailSender{&conf, &email, smtp.SendMail}
}

// Compatible with any SMTP server, either "Mailcather" or "SendGrid".
func (e *EmailSender) Send() error {
	// Password is required.
	if strings.TrimSpace(e.conf.Password) == "" {
		return errors.New("password is empty")
	}

	if e.email == nil {
		return errors.New("email is nil")
	}

	tpl, err := template.ParseFiles("templates/email.tpl")
	if err != nil {
		return err
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, e.email); err != nil {
		return err
	}

	text, err := ioutil.ReadAll(e.email.Reader)
	if err != nil {
		return err
	}

	msg := append(buff.Bytes(), "\r\n"...)
	msg = append(msg, text...)

	recipients := []string{e.email.ToAddress}

	auth := smtp.PlainAuth("", e.conf.Username, e.conf.Password, e.conf.Host)
	if err := e.send(e.conf.URL(), auth, e.email.FromAddress, recipients, msg); err != nil {
		return err
	}

	return nil
}
