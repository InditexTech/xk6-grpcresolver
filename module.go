package resolver

import (
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/resolver"
)

const (
	scheme                    = "k8s"
	grpcDebugLogsEnvVarName   = "GRPC_DEBUG_LOGS"
	grpcUpdateEveryEnvVarName = "GRCP_UPDATE_EVERY"
	trueValue                 = "true"
)

var updateEvery = 10 * time.Second
var showDebugLogs = false
var logger = logrus.New()

func init() {
	resolver.Register(&k8sGrpcModule{})
}

func logIfDebug(args ...interface{}) {
	if showDebugLogs {
		logger.Debug(args...)
	}
}
