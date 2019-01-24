package mailconsumer

import (
	"strings"

	"github.com/openware/postmaster/pkg/utils"
	"github.com/openware/postmaster/pkg/eventapi"
)

type tokenReceiverEvent struct {
	User  eventapi.User `json:"user"`
	Token string        `json:"token"`
}

// EmailConfirmationEvent is structure for processing "user.email.confirmation.token" event.
type EmailConfirmationEvent = tokenReceiverEvent

// ResetPasswordEvent is structure for processing "user.password.reset.token" event.
type ResetPasswordEvent = tokenReceiverEvent

// EmailConfirmationURI returns unique URL for user to confirm his identity.
func (event EmailConfirmationEvent) EmailConfirmationURI() string {
	url := utils.GetEnv("CONFIRM_URL", "http://example.com/#{}")
	return strings.Replace(url, "#{}", event.Token, 1)
}

// ResetPasswordURI returns unique URL for user to reset password.
func (event ResetPasswordEvent) ResetPasswordURI() string {
	url := utils.GetEnv("RESET_URL", "http://example.com/#{}")
	return strings.Replace(url, "#{}", event.Token, 1)
}
