package resolver

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/resolver"
)

type k8sGrpcModule struct{}

// Build parses the target for the service host and the endpoint port, returning an error if these can not be parsed.
// Should this succeed, it initialises a khsResolver, calls the first resolve and if this completes, it's returned.
func (kb *k8sGrpcModule) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	showDebugLogs = os.Getenv(grpcDebugLogsEnvVarName) == trueValue
	if os.Getenv(grpcUpdateEveryEnvVarName) != "" {
		updateEvery, _ = time.ParseDuration(os.Getenv(grpcUpdateEveryEnvVarName))
	}

	logIfDebug("Building GRPC resolver for target: ", target)
	strs := strings.Split(target.Endpoint(), ":")

	if len(strs) > 2 || len(strs) <= 0 {
		return nil, fmt.Errorf("couldn't parse given target endpoint: %s", target.Endpoint())
	}

	res := &k8sGrpcResolver{
		cc:          cc,
		serviceHost: strs[0],
		quitC:       make(chan struct{}),
	}

	if len(strs) == 2 {
		port, err := strconv.Atoi(strs[1])

		if err != nil {
			return nil, fmt.Errorf("couldn't parse given port: %s", strs[1])
		}

		res.endpointPort = port
	}

	err := res.resolve()

	if err != nil {
		return nil, err
	}

	go res.periodicUpdate()

	return res, nil
}

func (kb *k8sGrpcModule) Scheme() string {
	return scheme
}
