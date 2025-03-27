// SPDX-FileCopyrightText: 2025 INDUSTRIA DE DISEÑO TEXTIL S.A. (INDITEX S.A.)
//
// SPDX-License-Identifier: AGPL-3.0-only

package grpcresolver

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/resolver"
	"net"
	neturl "net/url"
	"testing"
	"time"
)

type BuilderTestSuite struct {
	suite.Suite
}

func TestBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}

func (suite *BuilderTestSuite) AfterTest(_, _ string) {
	hostsIPs = make(map[string][]net.IP)
}

func (suite *BuilderTestSuite) TestBuilderComplete() {
	host := "google.com"
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
	conn := TestClientConnImpl{
		stateUpdates: &Array[resolver.State]{},
	}

	resultResolver, err := builder.Build(target, conn, resolver.BuildOptions{})
	suite.Require().Nil(err)
	suite.Assert().Equal(conn, resultResolver.(*Resolver).conn)

	// Wait for resolver/updater to run
	time.Sleep(settings.UpdateEvery * 3)

	// Should have resolved the IPs in the singleton
	ips, ok := getResolverIPs(host)
	suite.Assert().True(ok)
	suite.Assert().NotEmpty(ips)

	// Should have updated the ClientConn with the resolved IPs
	connUpdates := conn.stateUpdates.values
	suite.Assert().GreaterOrEqual(len(connUpdates), 1)
	for _, connUpdate := range connUpdates {
		for _, address := range connUpdate.Addresses {
			suite.Assert().Equal(host, address.ServerName, fmt.Sprintf("asserting address.Servername %v", address))
		}
	}
}

func (suite *BuilderTestSuite) TestBuilderEmptyEndpoint() {
	builder := &Builder{}
	target := resolver.Target{}
	conn := TestClientConnImpl{}

	resultResolver, err := builder.Build(target, conn, resolver.BuildOptions{})
	suite.Assert().Nil(resultResolver)
	suite.Assert().Error(err)
	suite.Assert().ErrorContains(err, "invalid target")
}

func (suite *BuilderTestSuite) TestBuilderInvalidEndpoint() {
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
	suite.Assert().Nil(resultResolver)
	suite.Assert().Error(err)
	suite.Assert().ErrorContains(err, "invalid target endpoint \"google.com:aaa:50051\"")
}

func (suite *BuilderTestSuite) TestBuilderInvalidPort() {
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
	suite.Assert().Nil(resultResolver)
	suite.Assert().Error(err)
	suite.Assert().ErrorContains(err, "invalid port: aaa")
}
