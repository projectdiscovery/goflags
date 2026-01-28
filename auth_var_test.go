package goflags

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthVarWithValue(t *testing.T) {
	const envKey = "TEST_AUTH_KEY"
	os.Unsetenv(envKey)
	defer os.Unsetenv(envKey)

	var authValue string
	flagSet := NewFlagSet()
	flagSet.AuthVarP(&authValue, "auth", "a", envKey, "authentication token")

	err := flagSet.Parse("-auth=my-secret-token")
	assert.Nil(t, err)
	assert.Equal(t, "my-secret-token", authValue)
	assert.Equal(t, "my-secret-token", os.Getenv(envKey))
}

func TestAuthVarWithShortFlag(t *testing.T) {
	const envKey = "TEST_AUTH_KEY_SHORT"
	os.Unsetenv(envKey)
	defer os.Unsetenv(envKey)

	var authValue string
	flagSet := NewFlagSet()
	flagSet.AuthVarP(&authValue, "auth", "a", envKey, "authentication token")

	err := flagSet.Parse("-a=short-token")
	assert.Nil(t, err)
	assert.Equal(t, "short-token", authValue)
	assert.Equal(t, "short-token", os.Getenv(envKey))
}

func TestAuthVarFromExistingEnv(t *testing.T) {
	const envKey = "TEST_AUTH_KEY_ENV"
	os.Setenv(envKey, "env-token")
	defer os.Unsetenv(envKey)

	var authValue string
	flagSet := NewFlagSet()
	flagSet.AuthVarP(&authValue, "auth", "a", envKey, "authentication token")

	err := flagSet.Parse("")
	assert.Nil(t, err)
	assert.Equal(t, "env-token", authValue)
}

func TestAuthVarFlagOverridesEnv(t *testing.T) {
	const envKey = "TEST_AUTH_KEY_OVERRIDE"
	os.Setenv(envKey, "env-token")
	defer os.Unsetenv(envKey)

	var authValue string
	flagSet := NewFlagSet()
	flagSet.AuthVarP(&authValue, "auth", "a", envKey, "authentication token")

	err := flagSet.Parse("-auth=flag-token")
	assert.Nil(t, err)
	assert.Equal(t, "flag-token", authValue)
	assert.Equal(t, "flag-token", os.Getenv(envKey))
}
