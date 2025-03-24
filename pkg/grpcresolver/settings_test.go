// SPDX-FileCopyrightText: © 2025 Industria de Diseño Textil S.A. INDITEX
// SPDX-License-Identifier: Apache-2.0
package grpcresolver

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	changedEnvVars = make(map[string]string)
)

func TestLoadSettingsDefaults(t *testing.T) {
	err := LoadSettings()
	assert.Nil(t, err)
	assert.Equal(t, "k8s", settings.ProtocolName)
	assert.Equal(t, 3*time.Second, settings.UpdateEvery)
	assert.Equal(t, 3*time.Second, settings.SyncEvery)
	assert.False(t, settings.ShowDebugLogs)
}

func TestLoadSettingsFromEnv(t *testing.T) {
	defer clearEnvVars()

	resolverProtocol := "custom"
	setEnvVar("GRPC_RESOLVER_PROTOCOL", resolverProtocol)

	updateEvery := 65 * time.Second
	setEnvVar("GRPC_UPDATE_EVERY", "1m5s")

	syncEvery := 66 * time.Second
	setEnvVar("GRPC_SYNC_EVERY", "66s")

	setEnvVar("GRPC_DEBUG_LOGS", "true")

	err := LoadSettings()
	assert.Nil(t, err)
	assert.Equal(t, resolverProtocol, settings.ProtocolName)
	assert.Equal(t, updateEvery, settings.UpdateEvery)
	assert.Equal(t, syncEvery, settings.SyncEvery)
	assert.True(t, settings.ShowDebugLogs)
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
