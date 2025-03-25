// SPDX-FileCopyrightText: 2025 Industria de Diseño Textil S.A. INDITEX
//
// SPDX-License-Identifier: Apache-2.0

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
