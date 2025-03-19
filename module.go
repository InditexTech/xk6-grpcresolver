package resolver

import (
	"github.com/InditexTech/xk6-grpcresolver/pkg/grpcresolver"
	"google.golang.org/grpc/resolver"
)

// Register the custom gRPC resolver.
func init() {
	if err := grpcresolver.LoadSettings(); err != nil {
		panic("failed to load settings: " + err.Error())
	}

	resolver.Register(&grpcresolver.Builder{})
}
