package grpcresolver

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func logIfDebug(args ...interface{}) {
	if settings.ShowDebugLogs {
		logger.Info(args...)
	}
}
