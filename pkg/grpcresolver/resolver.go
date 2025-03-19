package grpcresolver

import (
	"fmt"
	"google.golang.org/grpc/resolver"
	"net"
	"sync"
	"time"
)

var (
	// hostsIPs contains an updated list of the IPs resolved per host.
	// The existence of a host in this map means the periodic resolver is running for that host.
	hostsIPs                  = make(map[string][]net.IP)
	hostsIPsLock              sync.Mutex
	periodicResolverStartLock sync.Mutex
)

// Resolver implements the gRPC client resolver.go Resolver interface, so can replace the default implementation in the k6 gRPC client.
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

// update updates the Resolver addresses with the current resolverIps list.
func (r *Resolver) update() error {
	newIps := r.containsNewIp()
	resolverIps, _ := getResolverIPs(r.serviceHost)
	deletedIps := len(r.currentIps) > len(resolverIps)
	same := !newIps && !deletedIps

	if same {
		logIfDebug(fmt.Sprintf("No changes in resolved IPs for %s. Current IPs: %v", r.serviceHost, r.currentIps))
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
	resolverIps, _ := getResolverIPs(r.serviceHost)
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
	logIfDebug("Starting periodic updater for ", r.serviceHost)
	go r.periodicUpdaterTask()
}

func (r *Resolver) periodicUpdaterTask() {
	ticker := time.NewTicker(settings.SyncEvery)
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
	periodicResolverStartLock.Lock()
	defer periodicResolverStartLock.Unlock()

	if _, resolverStarted := getResolverIPs(serviceHost); resolverStarted {
		return
	}

	logIfDebug("Starting periodic resolver for ", serviceHost)
	setResolverIPs(serviceHost, make([]net.IP, 0))
	go func() {
		periodicResolverTask(serviceHost)
		for range time.NewTicker(settings.UpdateEvery).C {
			periodicResolverTask(serviceHost)
		}
	}()
}

func periodicResolverTask(serviceHost string) {
	ips, err := net.LookupIP(serviceHost)
	if err != nil {
		logger.Error(fmt.Sprintf("Error looking up IPs for %s: %s", serviceHost, err.Error()))
	} else {
		logIfDebug(fmt.Sprintf("Looking up IPs for %s: %s", serviceHost, ips))
		setResolverIPs(serviceHost, ips)
	}
}

func getResolverIPs(serviceHost string) ([]net.IP, bool) {
	ips, ok := hostsIPs[serviceHost]
	return ips, ok
}

func setResolverIPs(serviceHost string, ips []net.IP) {
	hostsIPsLock.Lock()
	defer hostsIPsLock.Unlock()

	hostsIPs[serviceHost] = ips
}
