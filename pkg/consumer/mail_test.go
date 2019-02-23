package consumer

import (
	"net/smtp"
	"testing"

	"github.com/stretchr/testify/assert"
)

type emailRecorder struct {
	addr string
	auth smtp.Auth
	from string
	to   []string
	msg  []byte
}

func mockSend(errToReturn error) (func(string, smtp.Auth, string, []string, []byte) error, *emailRecorder) {
	req := new(emailRecorder)
	return func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		*req = emailRecorder{addr, a, from, to, msg}
		return errToReturn
	}, req
}

func TestEmailSender_Send(t *testing.T) {
	fakeEmail := Email{
		FromAddress: "test@postmaster.com",
		FromName:    "Postmaster",
		ToAddress:   "johndoe@gmail.com",
		Subject:     "Test",
		Reader:      nil,
	}

	t.Run("returns error, if password is empty", func(t *testing.T) {
		f, _ := mockSend(nil)
		sender := &EmailSender{send: f, email: &fakeEmail, conf: &SMTPConf{Password: ""}}
		err := sender.Send()

		assert.Equal(t, "password is empty", err.Error())
	})
}
