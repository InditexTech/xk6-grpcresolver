package grpcresolver

import (
	"fmt"
	"google.golang.org/grpc/resolver"
	"net"
	"strconv"
	"strings"
	"time"
)

// Builder implements the gRPC resolver.go Builder interface, so can replace the default implementation in the k6 gRPC client.
type Builder struct{}

// Build main logic entrypoint for the plugin. Is called when a VU calls client.connect, and creates the Resolver used by the k6 gRPC client.
func (b *Builder) Build(target resolver.Target, conn resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	if err := loadSettings(); err != nil {
		return nil, fmt.Errorf("failed loading settings: %w", err)
	}

	logger.Info("Building GRPC resolver for target: ", target)

	endpointChunks := strings.Split(target.Endpoint(), ":")
	if len(endpointChunks) > 2 || len(endpointChunks) <= 0 {
		return nil, fmt.Errorf("couldn't parse given target endpoint: %s", target.Endpoint())
	}

	endpointHost := endpointChunks[0]
	customResolver := &Resolver{
		conn:        conn,
		serviceHost: endpointHost,
		quitC:       make(chan struct{}),
	}

	if len(endpointChunks) == 2 {
		portStr := endpointChunks[1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse given port: %s", portStr)
		}

		customResolver.endpointPort = port
	}

	if err := customResolver.update(); err != nil {
		return nil, err
	}

	startPeriodicResolver(endpointHost)
	customResolver.startPeriodicUpdater()

	return customResolver, nil
}

// Scheme returns the configured ProtocolName.
func (b *Builder) Scheme() string {
	return settings.ProtocolName
}

func (b *Builder) periodicResolveServiceIps(serviceHost string) {
	t := time.NewTicker(10 * time.Second)
	for range t.C {
		if err := b.resolveServiceIps(serviceHost); err != nil {
			logger.Error("Error looking up IPs for ", serviceHost, " : ", err)
		}
	}
}

func (b *Builder) resolveServiceIps(serviceHost string) error {
	var err error
	resolverIps, err = net.LookupIP(serviceHost)
	logIfDebug("Looking up IPs for ", serviceHost, " : ", resolverIps)

	return err
}
