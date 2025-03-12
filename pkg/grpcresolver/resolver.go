package grpcresolver

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/resolver"
)

const (
	name                      = "k8s"
	grpcDebugLogsEnvVarName   = "GRPC_DEBUG_LOGS"
	grpcUpdateEveryEnvVarName = "GRCP_UPDATE_EVERY"
	trueValue                 = "true"
)

var updateEvery = 3 * time.Second
var showDebugLogs = false
var logger = logrus.New()
var isResolverStarted = false
var resolverIps []net.IP

// Register k8s with gRPC.
func init() {
	resolver.Register(&K8sBuilder{})
}

type K8sBuilder struct{}

// Build parses the target for the service host and the endpoint port, returning an error if these can not be parsed.
// Should this succeed, it initialises a khsResolver, calls the first resolve and if this completes, it's returned.
func (kb *K8sBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	showDebugLogs = os.Getenv(grpcDebugLogsEnvVarName) == trueValue
	if os.Getenv(grpcUpdateEveryEnvVarName) != "" {
		updateEvery, _ = time.ParseDuration(os.Getenv(grpcUpdateEveryEnvVarName))
	}

	logger.Info("Building GRPC resolver for target: ", target)
	strs := strings.Split(target.Endpoint(), ":")

	if len(strs) > 2 || len(strs) <= 0 {
		return nil, fmt.Errorf("couldn't parse given target endpoint: %s", target.Endpoint())
	}

	updater := &k8sClientUpdater{
		cc:          cc,
		serviceHost: strs[0],
		quitC:       make(chan struct{}),
	}

	if len(strs) == 2 {
		port, err := strconv.Atoi(strs[1])

		if err != nil {
			return nil, fmt.Errorf("couldn't parse given port: %s", strs[1])
		}

		updater.endpointPort = port
	}

	err := updater.update()

	if err != nil {
		return nil, err
	}

	if !isResolverStarted {
		kr := &k8sResolver{}
		isResolverStarted = true
		kr.resolveServiceIps(strs[0])
		go kr.periodicResolveServiceIps(strs[0])
	}

	go updater.periodicUpdateClient()

	return updater, nil
}

// Scheme returns `hpa`.
func (kb *K8sBuilder) Scheme() string {
	return name
}

type k8sResolver struct {
}

func (k *k8sResolver) periodicResolveServiceIps(serviceHost string) {
	t := time.NewTicker(10 * time.Second)
	for range t.C {
		if err := k.resolveServiceIps(serviceHost); err != nil {
			logger.Error("Error looking up IPs for ", serviceHost, " : ", err)
		}
	}
}

func (k *k8sResolver) resolveServiceIps(serviceHost string) error {
	var err error
	resolverIps, err = net.LookupIP(serviceHost)
	logIfDebug("Looking up IPs for ", serviceHost, " : ", resolverIps)

	return err
}

// k8sClientUpdater is the resolver for Kubernetes Headless Services. When called, it looks up all the A Records for the given
// Host, passing them to a ClientConn as Backends.
type k8sClientUpdater struct {
	cc resolver.ClientConn

	serviceHost  string
	endpointPort int
	currentIps   []net.IP

	quitC chan struct{}
}

// update calls the ClientConn NewAddress callback with all the IPs returned from a standard DNS lookup to the serviceHost,
// affixing kr.endpointPort to each of them.
func (kr *k8sClientUpdater) update() error {
	newIps := kr.containsNewIp()
	deletedIps := len(kr.currentIps) > len(resolverIps)
	same := !newIps && !deletedIps

	if same {
		logIfDebug("No changes in resolved IPs")
		return nil
	}

	addrs := make([]resolver.Address, len(resolverIps))

	for i, ip := range resolverIps {
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
	if len(kr.currentIps) > 0 {
		logger.Info("Service host k8s:///", kr.serviceHost, " has been resolved successfully with IPs ", addrs)
	}
	kr.currentIps = resolverIps
	kr.cc.UpdateState(resolver.State{
		Addresses:     addrs,
		ServiceConfig: kr.cc.ParseServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	})

	return nil
}

func (kr *k8sClientUpdater) containsNewIp() bool {
	newIps := false
	for _, ip := range resolverIps {
		exists := false
		for _, currentIp := range kr.currentIps {
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

// periodicUpdateClient periodically calls resolve to ensure kr.cc contains an recent list of the service endpoints.
func (kr *k8sClientUpdater) periodicUpdateClient() {
	t := time.NewTicker(updateEvery)
	for {
		select {
		case <-t.C:
			err := kr.update()

			if err != nil {
				logger.Error("[hpa - resolver.go] error resolving: ", err)
			}
		case <-kr.quitC:
			return
		}
	}
}

func (kr *k8sClientUpdater) ResolveNow(option resolver.ResolveNowOptions) {
	err := kr.update()

	if err != nil {
		logger.Error("[hpa - resolver.go] error resolving: ", err)
	}
}

func (kr *k8sClientUpdater) Close() {
	kr.quitC <- struct{}{}
}

func logIfDebug(args ...interface{}) {
	if showDebugLogs {
		logger.Info(args...)
	}
}
