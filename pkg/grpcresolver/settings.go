package grpcresolver

import (
	"github.com/caarlos0/env/v11"
	"time"
)

var settings SettingsSpec

type SettingsSpec struct {
	ProtocolName  string        `env:"GRPC_RESOLVER_PROTOCOL" envDefault:"k8s"`
	UpdateEvery   time.Duration `env:"GRPC_UPDATE_EVERY" envDefault:"3s"`
	SyncEvery     time.Duration `env:"GRPC_SYNC_EVERY" envDefault:"3s"`
	ShowDebugLogs bool          `env:"GRPC_DEBUG_LOGS" envDefault:"false"`
}

func loadSettings() error {
	return env.Parse(&settings)
}
