package consumer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmailConfirmationEvent_EmailConfirmationURI(t *testing.T) {
	res := EmailConfirmationEvent{
		Token: "ixj717iex",
	}.EmailConfirmationURI()

	assert.Contains(t, res, "ixj717iex")
}

func TestPasswordResetEvent_ResetPasswordURI(t *testing.T) {
	res := EmailConfirmationEvent{
		Token: "yxy1U1yxy",
	}.ResetPasswordURI()

	assert.Contains(t, res, "yxy1U1yxy")
}
