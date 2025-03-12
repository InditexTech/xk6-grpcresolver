package grpcresolver

import (
	"github.com/caarlos0/env/v11"
)

var settings SettingsSpec

type SettingsSpec struct {
	ProtocolName  string `env:"GRPC_RESOLVER_PROTOCOL" envDefault:"k8s"`
	ShowDebugLogs bool
}

func loadSettings() error {
	return env.Parse(&settings)
}
