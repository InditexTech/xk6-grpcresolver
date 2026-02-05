// SPDX-FileCopyrightText: 2025 INDUSTRIA DE DISEÑO TEXTIL S.A. (INDITEX S.A.)
//
// SPDX-License-Identifier: AGPL-3.0-only

package grpcresolver

import (
	"time"

	"github.com/mstoykov/envconfig"
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
