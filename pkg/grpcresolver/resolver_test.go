// SPDX-FileCopyrightText: 2025 INDUSTRIA DE DISEÑO TEXTIL S.A. (INDITEX S.A.)
//
// SPDX-License-Identifier: AGPL-3.0-only

package grpcresolver

import (
	"fmt"
	"math/rand/v2"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PeriodicResolverTestSuite struct {
	suite.Suite
}

func TestPeriodicResolverTestSuite(t *testing.T) {
	suite.Run(t, new(PeriodicResolverTestSuite))
}

func (suite *PeriodicResolverTestSuite) AfterTest(_, _ string) {
	hostsIPs.Clear()
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

// Read and Write the hostsIPs map several times, concurrently from two goroutines.
// Should not panic because of the concurrent access to the map.
func (suite *PeriodicResolverTestSuite) TestGetSetConcurrently() {
	writeTimes := 1000000
	getTimes := 1000000
	randHostnamesCount := 20

	getRandHostname := func() string {
		n := rand.IntN(randHostnamesCount)
		return fmt.Sprintf("%d.test.local", n)
	}

	getRandIPs := func() []net.IP {
		ips := make([]net.IP, 5)
		for i := 0; i < 5; i++ {
			ip := net.IPv4(byte(rand.IntN(256)), byte(rand.IntN(256)), byte(rand.IntN(256)), byte(rand.IntN(256)))
			ips = append(ips, ip)
		}
		return ips
	}

	wait := sync.WaitGroup{}
	wait.Add(2)

	// Write
	go func() {
		for i := 0; i < writeTimes; i++ {
			setResolverIPs(getRandHostname(), getRandIPs())
		}
		wait.Done()
	}()

	// Read
	go func() {
		for i := 0; i < getTimes; i++ {
			getResolverIPs(getRandHostname())
		}
		wait.Done()
	}()

	wait.Wait()
}
