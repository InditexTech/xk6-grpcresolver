// SPDX-FileCopyrightText: © 2025 Industria de Diseño Textil S.A. INDITEX
// SPDX-License-Identifier: Apache-2.0
package grpcresolver

import (
	"github.com/mstoykov/envconfig"
	"time"
)

var settings = &SettingsSpec{}

type SettingsSpec struct {
	ProtocolName  string        `envconfig:"GRPC_RESOLVER_PROTOCOL" default:"k8s"`
	UpdateEvery   time.Duration `envconfig:"GRPC_UPDATE_EVERY" default:"3s"`
	SyncEvery     time.Duration `envconfig:"GRPC_SYNC_EVERY" default:"3s"`
	ShowDebugLogs bool          `envconfig:"GRPC_DEBUG_LOGS" default:"false"`
}

func LoadSettings() error {
	return envconfig.Process("", settings)
}
