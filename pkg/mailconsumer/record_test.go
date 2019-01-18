package mailconsumer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountRecordConfirmationUri(t *testing.T) {
	res := AccountCreatedEvent{
		Token: "12345",
	}.ConfirmationURI()

	assert.Contains(t, res, "12345")
}
