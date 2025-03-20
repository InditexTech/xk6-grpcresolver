package grpcresolver

import (
	"github.com/stretchr/testify/suite"
	"net"
	"testing"
)

type PeriodicResolverTestSuite struct {
	suite.Suite
}

func TestPeriodicResolverTestSuite(t *testing.T) {
	suite.Run(t, new(PeriodicResolverTestSuite))
}

func (suite *PeriodicResolverTestSuite) AfterTest(_, _ string) {
	hostsIPs = make(map[string][]net.IP)
}

func (suite *PeriodicResolverTestSuite) TestGetSetResolverIPs() {
	hostname := "foo.bar"
	ips := []net.IP{net.ParseIP("8.8.8.8"), net.ParseIP("8.8.4.4")}
	setResolverIPs(hostname, ips)

	resultIps, ok := getResolverIPs(hostname)
	suite.Assert().True(ok)
	suite.Assert().Equal(ips, resultIps)
}

func (suite *PeriodicResolverTestSuite) TestGetResolverIPsNotRegistered() {
	ips, ok := getResolverIPs("foo.bar")
	suite.Assert().Empty(ips)
	suite.Assert().False(ok)
}

func (suite *PeriodicResolverTestSuite) TestGetResolverIPsEmpty() {
	hostname := "foo.bar"
	setResolverIPs(hostname, make([]net.IP, 0))

	ips, ok := getResolverIPs(hostname)
	suite.Assert().Empty(ips)
	suite.Assert().True(ok)
}

func (suite *PeriodicResolverTestSuite) TestPeriodicResolverTask() {
	hostname := "google.com"
	ips, ok := getResolverIPs(hostname)
	suite.Require().Empty(ips)
	suite.Require().False(ok)

	runLookupTaskOnce(hostname)

	ips, ok = getResolverIPs(hostname)
	suite.Require().True(ok)
	suite.Require().GreaterOrEqual(len(ips), 1)
}
