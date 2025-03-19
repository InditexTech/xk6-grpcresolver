package grpcresolver

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func logIfDebug(args ...interface{}) {
	if settings.ShowDebugLogs {
		// TODO: show as Debug in k6. Not showing when using logger.Debug & running k6 with -v.
		logger.Info(args...)
	}
}
