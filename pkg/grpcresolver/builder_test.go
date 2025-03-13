package grpcresolver

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/resolver"
	neturl "net/url"
	"testing"
	"time"
)

// TestClientConnImpl struct that implements resolver.ClientConn
// Copied from google.golang.org/grpc@v1.60.0/internal/testutils/resolver.go
type TestClientConnImpl struct {
	resolver.ClientConn
}

func TestCompleteBuilder(t *testing.T) {
	endpoint := "k8s:///google.com:50051"
	url, err := neturl.Parse(endpoint)
	if err != nil {
		panic(err)
	}

	if err := LoadSettings(); err != nil {
		panic(err)
	}

	settings.UpdateEvery = 100 * time.Millisecond
	settings.SyncEvery = settings.UpdateEvery

	builder := &Builder{}
	target := resolver.Target{URL: *url}
	conn := TestClientConnImpl{}

	resultResolver, err := builder.Build(target, conn, resolver.BuildOptions{})
	require.Nil(t, err)
	assert.Equal(t, conn, resultResolver.(*Resolver).conn)

	// TODO Try to continue the test

	//// Wait for resolver/updater to run
	//time.Sleep(settings.UpdateEvery * 3)
	//
	//ips, ok := getResolverIPs("google.com")
	//assert.True(t, ok)
	//assert.NotEmpty(t, ips)
}

func TestBuilderEmptyEndpoint(t *testing.T) {
	builder := &Builder{}
	target := resolver.Target{}
	conn := TestClientConnImpl{}

	resultResolver, err := builder.Build(target, conn, resolver.BuildOptions{})
	assert.Nil(t, resultResolver)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid target")
}

func TestBuilderInvalidEndpoint(t *testing.T) {
	endpoint := "k8s:///google.com:aaa:50051"
	url, err := neturl.Parse(endpoint)
	if err != nil {
		panic(err)
	}

	if err := LoadSettings(); err != nil {
		panic(err)
	}

	builder := &Builder{}
	target := resolver.Target{URL: *url}
	conn := TestClientConnImpl{}

	resultResolver, err := builder.Build(target, conn, resolver.BuildOptions{})
	assert.Nil(t, resultResolver)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid target endpoint \"google.com:aaa:50051\"")
}

func TestBuilderInvalidPort(t *testing.T) {
	endpoint := "k8s:///google.com:aaa"
	url, err := neturl.Parse(endpoint)
	if err != nil {
		panic(err)
	}

	if err := LoadSettings(); err != nil {
		panic(err)
	}

	builder := &Builder{}
	target := resolver.Target{URL: *url}
	conn := TestClientConnImpl{}

	resultResolver, err := builder.Build(target, conn, resolver.BuildOptions{})
	assert.Nil(t, resultResolver)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid port: aaa")
}
