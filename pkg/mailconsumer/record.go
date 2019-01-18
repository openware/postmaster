package mailconsumer

import "strings"

import (
	"github.com/shal/pigeon/pkg/eventapi"
	"github.com/shal/pigeon/pkg/utils"
)

type AccountCreatedEvent struct {
	User  eventapi.User `json:"user"`
	Token string        `json:"token"`
}

func (r AccountCreatedEvent) ConfirmationURI() string {
	url := utils.GetEnv("CONFIRM_URL",
		"http://www.example.com/accounts/confirmation?confirmation_token=#{}",
	)

	return strings.Replace(url, "#{}", r.Token, 1)
}
