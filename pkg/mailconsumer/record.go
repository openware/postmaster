package mailconsumer

import (
	"github.com/shal/pigeon/pkg/eventapi"
	"strings"
)

import (
	"github.com/shal/pigeon/pkg/utils"
)

type AccountCreatedEvent struct {
	User  eventapi.User `json:"user"`
	Token string        `json:"string"`
}

func (r AccountCreatedEvent) ConfirmationURI() string {
	url := utils.GetEnv("CONFIRM_URL",
		"http://www.example.com/accounts/confirmation?confirmation_token=#{}",
	)

	return strings.Replace(url, "#{}", r.Token, 1)
}
