// SPDX-FileCopyrightText: 2025 INDUSTRIA DE DISEÑO TEXTIL S.A. (INDITEX S.A.)
//
// SPDX-License-Identifier: AGPL-3.0-only

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
