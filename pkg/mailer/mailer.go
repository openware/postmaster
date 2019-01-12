package mailer

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
)

type MailTpl struct {
	To   string
	From string
}

func randomString(length int) string {
	buf := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

func (m MailTpl) MessageID() string {
	host, err := os.Hostname()

	if err != nil {
		log.Println()
	}

	return fmt.Sprintf("%s@%s", randomString(32), host)
}
