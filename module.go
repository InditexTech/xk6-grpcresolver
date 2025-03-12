package resolver

import (
	"github.com/InditexTech/xk6-grpcresolver/pkg/grpcresolver"
	"google.golang.org/grpc/resolver"
)

// Register k8s with gRPC.
func init() {
	resolver.Register(&grpcresolver.Builder{})
}
