package resolver

import (
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc/resolver"
)

// k8sGrpcResolver is the resolver for Kubernetes Headless Services. When called, it looks up all the A Records for the given
// Host, passing them to a ClientConn as Backends.
type k8sGrpcResolver struct {
	cc resolver.ClientConn

	serviceHost  string
	endpointPort int
	currentIps   []net.IP

	quitC chan struct{}
}

// resolve calls the ClientConn NewAddress callback with all the IPs returned from a standard DNS lookup to the serviceHost,
// affixing kr.endpointPort to each of them.
func (kr *k8sGrpcResolver) resolve() error {
	ips, err := net.LookupIP(kr.serviceHost)

	if err != nil {
		return err
	}

	logIfDebug("Resolved IPs: ", ips, " Current IPs: ", kr.currentIps)

	same := true
	for _, ip := range ips {
		exists := false
		for _, currentIp := range kr.currentIps {
			if ip.Equal(currentIp) {
				exists = true
				break
			}
		}
		if !exists {
			logIfDebug("New IP found: ", ip)
			same = false
			break
		}
	}

	if same {
		logIfDebug("No changes in resolved IPs")
		return nil
	}

	kr.currentIps = ips
	addrs := make([]resolver.Address, len(ips))

	for i, ip := range ips {
		addr := ip.String()

		if kr.endpointPort != 0 {
			addr = fmt.Sprintf("%s:%d", addr, kr.endpointPort)
		}

		addrs[i] = resolver.Address{
			Addr:       addr,
			ServerName: kr.serviceHost,
		}
	}

	// NOTE: Use of the built-in Round Robin Balancer (google.golang.org/grpc/balancer/roundrobin) is now set via
	// ServiceConfig JSON instead of the depreciated grpc.WithBalancerName(roundrobin.Name), previously a client DialOption.
	// However, the gRPC Service Config docs (https://github.com/grpc/grpc/blob/master/doc/service_config.md) suggest
	// loadBalancingPolicy is also being deprecated with no clear alternative.
	//
	// grpc/service_config.go currently supports a 'loadBalancingConfig' field, however it looks likely to change, so for
	// now stick to the existing JSON definition.
	logIfDebug("Service host k8s:///", kr.serviceHost, " has been resolved successfully with IPs ", addrs)
	kr.cc.UpdateState(resolver.State{
		Addresses:     addrs,
		ServiceConfig: kr.cc.ParseServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	})

	return nil
}

// periodicUpdate periodically calls resolve to ensure kr.cc contains an recent list of the service endpoints.
func (kr *k8sGrpcResolver) periodicUpdate() {
	t := time.NewTicker(updateEvery)
	for {
		select {
		case <-t.C:
			err := kr.resolve()

			if err != nil {
				logger.Error("[hpa - resolver.go] error resolving: ", err)
			}
		case <-kr.quitC:
			return
		}
	}
}

// Resolve now runs an internal resolve, updating hpa.cc with the current list of endpoints.
func (kr *k8sGrpcResolver) ResolveNow(option resolver.ResolveNowOptions) {
	err := kr.resolve()

	if err != nil {
		logger.Error("[hpa - resolver.go] error resolving: ", err)
	}
}

func (kr *k8sGrpcResolver) Close() {
	kr.quitC <- struct{}{}
}
