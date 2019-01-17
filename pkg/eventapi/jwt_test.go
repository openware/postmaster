package eventapi

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestParseJWT(t *testing.T) {
	FakeKeyFunc := func(_ *jwt.Token) (interface{}, error) {
		pub, err := ioutil.ReadFile("../../test/sample.key.pub")
		assert.NoError(t, err)
		signingKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
		assert.NoError(t, err)
		return signingKey, nil
	}

	t.Run("Invalid token", func(t *testing.T) {
		claims, err := ParseJWT("x.y", FakeKeyFunc)
		assert.Equal(t, "token contains an invalid number of segments", err.Error())
		assert.Nil(t, claims)
	})

	t.Run("Valid token", func(t *testing.T) {
		token := "eyJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJiYXJvbmciLCJqdGkiOiIxMTFhYTExYS0xMTFhLTExMTEtYWExMS0xMWFhMTExYTExYTEiLCJpYXQiOjE1NDc3MjcwMjksImV2ZW50Ijp7InJlY29yZCI6eyJ1aWQiOiJJREpPSE5ET0UxMjMiLCJlbWFpbCI6ImpvaG5AZG9lLmNvbSIsImxldmVsIjowLCJvdHBfZW5hYmxlZCI6ZmFsc2UsImNvbmZpcm1hdGlvbl90b2tlbiI6IjEyMzQ1IiwiY29uZmlybWF0aW9uX3NlbnRfYXQiOiIyMDE5LTAxLTE3VDEyOjEwOjI5WiIsInN0YXRlIjoicGVuZGluZyIsImZhaWxlZF9hdHRlbXB0cyI6MCwiY3JlYXRlZF9hdCI6IjIwMTktMDEtMTdUMTI6MTA6MjlaIiwidXBkYXRlZF9hdCI6IjIwMTktMDEtMTdUMTI6MTA6MjlaIn0sIm5hbWUiOiJtb2RlbC5hY2NvdW50LmNyZWF0ZWQifX0.Ci4-8Af4kd3QAhxe2eLL1zbakz208NjRiAJ6iiQXEi2L9izIeFn90lu7Pn0QiRh1O_FdJFRcyRvUVdMpNjNZCGVJsjQWE_lmT1nIoR9AMNK5nsHBzP2Ibs28nEgZZntPo_F0Z4F-k4FGbCjcNeF76szLtrHsSy7tcmrKALrsqKGrwh4IcE24VAbHQBTLs2nFOhZPdaQzcZlV8ExUBGVd9oHQDNffPIFLM_U3TjUIQFKWovKQ-gGPPTDlqzbIJgrp7xflHYh_lMKYU0_ZSQRMXRZyIcwbKVea0Jc6GC-daEXJ8PAgXDxGx4FLeqLIQt62qa9Ysifd6onOnFHjw8SK2g"
		claims, err := ParseJWT(token, FakeKeyFunc)

		assert.NoError(t, err)
		assert.Nil(t, claims.Event.Changes)
		assert.NotNil(t, claims.Event.Record)
		assert.Equal(t, "model.account.created", claims.Event.Name)
	})
}
