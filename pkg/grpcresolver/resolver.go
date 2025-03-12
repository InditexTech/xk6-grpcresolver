package grpcresolver

import (
	"fmt"
	"google.golang.org/grpc/resolver"
	"net"
	"sync"
	"time"
)

var (
	// resolverIps contains an updated list of the IPs resolved for the serviceHost.
	resolverIps []net.IP

	periodicResolverStarted     bool
	periodicResolverStartedLock sync.Mutex
)

// Resolver implements the gRPC client resolver.go Resolved interface, so can replace the default implementation in the k6 gRPC client.
type Resolver struct {
	conn resolver.ClientConn

	serviceHost  string
	endpointPort int
	currentIps   []net.IP

	quitC chan struct{}
}

// ResolveNow runs an internal resolve, updating with the current list of endpoints.
func (r *Resolver) ResolveNow(_ resolver.ResolveNowOptions) {
	if err := r.update(); err != nil {
		logger.Error("error resolving: ", err)
	}
}

func (r *Resolver) Close() {
	r.quitC <- struct{}{}
}

// update updates the Resolver addressed with the current resolverIps list.
func (r *Resolver) update() error {
	newIps := r.containsNewIp()
	deletedIps := len(r.currentIps) > len(resolverIps)
	same := !newIps && !deletedIps

	if same {
		logIfDebug("No changes in resolved IPs")
		return nil
	}

	addrs := make([]resolver.Address, len(resolverIps))

	for i, ip := range resolverIps {
		addr := ip.String()

		if r.endpointPort != 0 {
			addr = fmt.Sprintf("%s:%d", addr, r.endpointPort)
		}

		addrs[i] = resolver.Address{
			Addr:       addr,
			ServerName: r.serviceHost,
		}
	}

	// NOTE: Use of the built-in Round Robin Balancer (google.golang.org/grpc/balancer/roundrobin) is now set via
	// ServiceConfig JSON instead of the depreciated grpc.WithBalancerName(roundrobin.Name), previously a client DialOption.
	// However, the gRPC Service Config docs (https://github.com/grpc/grpc/blob/master/doc/service_config.md) suggest
	// loadBalancingPolicy is also being deprecated with no clear alternative.
	//
	// grpc/service_config.go currently supports a 'loadBalancingConfig' field, however it looks likely to change, so for
	// now stick to the existing JSON definition.
	if len(r.currentIps) > 0 {
		logger.Info("Service host k8s:///", r.serviceHost, " has been resolved successfully with IPs ", addrs)
	}
	r.currentIps = resolverIps
	_ = r.conn.UpdateState(resolver.State{
		Addresses:     addrs,
		ServiceConfig: r.conn.ParseServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	})

	return nil
}

func (r *Resolver) containsNewIp() bool {
	newIps := false
	for _, ip := range resolverIps {
		exists := false
		for _, currentIp := range r.currentIps {
			if ip.Equal(currentIp) {
				exists = true
				break
			}
		}
		if !exists {
			logIfDebug("New IP found: ", ip)
			newIps = true
			break
		}
	}
	return newIps
}

// startPeriodicUpdater starts a task that periodically synchronizes the IPs of the Resolver with those in the resolverIps array.
func (r *Resolver) startPeriodicUpdater() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			if err := r.update(); err != nil {
				logger.Error("periodic updater failed resolving: ", err)
			}
		case <-r.quitC:
			return
		}
	}
}

// startPeriodicResolver starts a task that periodically analyzes the IPs of the serviceHost.
// The task is initialized only once for all VUs.
// The IPs are stored in resolverIps singleton.
func startPeriodicResolver(serviceHost string) {
	periodicResolverStartedLock.Lock()
	defer periodicResolverStartedLock.Unlock()

	if periodicResolverStarted {
		return
	}

	go func() {
		// TODO configurable period time
		// TODO Check if the ticker runs for the first time then waits, or viceversa
		for range time.NewTicker(10 * time.Second).C {
			periodicResolverTask(serviceHost)
		}
	}()
	periodicResolverStarted = true
}

func periodicResolverTask(serviceHost string) {
	var err error
	resolverIps, err = net.LookupIP(serviceHost)
	if err != nil {
		logger.Error("Error looking up IPs for ", serviceHost, " : ", err)
	} else {
		logIfDebug("Looking up IPs for ", serviceHost, " : ", resolverIps)
	}
}
