package mailconsumer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountRecordConfirmationUri(t *testing.T) {
	res := AccountRecord{
		ConfirmationToken: "12345",
	}.ConfirmationUri()

	assert.Contains(t, res, "12345")
}
