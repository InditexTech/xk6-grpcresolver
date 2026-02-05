// SPDX-FileCopyrightText: 2025 INDUSTRIA DE DISEÑO TEXTIL S.A. (INDITEX S.A.)
//
// SPDX-License-Identifier: AGPL-3.0-only

package grpcresolver

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/grpc/resolver"
)

// Builder implements the gRPC resolver.go Builder interface, so can replace the default implementation in the k6 gRPC client.
type Builder struct{}

// Build main logic entrypoint for the plugin. Is called when a VU calls client.connect, and creates the Resolver used by the k6 gRPC client.
func (b *Builder) Build(target resolver.Target, conn resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	endpoint := target.Endpoint()
	if endpoint == "" {
		return nil, fmt.Errorf("invalid target \"%s\"", target.String())
	}

	endpointChunks := strings.Split(endpoint, ":")
	if len(endpointChunks) > 2 || len(endpointChunks) <= 0 {
		return nil, fmt.Errorf("invalid target endpoint \"%s\"", endpoint)
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
			return nil, fmt.Errorf("invalid port: %s", portStr)
		}

		customResolver.endpointPort = port
	}

	if err := customResolver.update(); err != nil {
		return nil, err
	}

	startPeriodicLookupTask(endpointHost)
	customResolver.startPeriodicSyncTask()

	return customResolver, nil
}

// Scheme returns the configured ProtocolName.
func (b *Builder) Scheme() string {
	return settings.ProtocolName
}
