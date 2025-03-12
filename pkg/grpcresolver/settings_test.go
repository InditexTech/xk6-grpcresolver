package grpcresolver

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	changedEnvVars = make(map[string]string)
)

func TestLoadSettingsDefaults(t *testing.T) {
	err := loadSettings()
	assert.Nil(t, err)
	assert.Equal(t, "k8s", settings.ProtocolName)
}

func TestLoadSettingsFromEnv(t *testing.T) {
	defer clearEnvVars()

	resolverProtocol := "custom"
	setEnvVar("GRPC_RESOLVER_PROTOCOL", resolverProtocol)

	err := loadSettings()
	assert.Nil(t, err)
	assert.Equal(t, resolverProtocol, settings.ProtocolName)
}

func setEnvVar(key string, value string) {
	if err := os.Setenv(key, value); err != nil {
		panic(err)
	}
	changedEnvVars[key] = value
}

func clearEnvVars() {
	for key := range changedEnvVars {
		if err := os.Unsetenv(key); err != nil {
			panic(err)
		}
		delete(changedEnvVars, key)
	}
}
