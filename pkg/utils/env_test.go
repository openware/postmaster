package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Should panic.
func TestPanicMustGetEnv(t *testing.T) {
	expected := "Environment variable INVALID_ENV does not set"
	assert.PanicsWithValue(t, expected, func() {
		MustGetEnv("INVALID_ENV")
	})
}

// Should return fallback.
func TestInvalidEnvGetEnv(t *testing.T) {
	expected := "default"
	res := GetEnv("INVALID_ENV", expected)

	assert.Equal(t, expected, res)
}
