package grpcresolver

import (
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func logIfDebug(args ...interface{}) {
	if settings.ShowDebugLogs {
		Logger.Info(args...)
	}
}
