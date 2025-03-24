// SPDX-FileCopyrightText: © 2025 Industria de Diseño Textil S.A. INDITEX
// SPDX-License-Identifier: APACHE-2.0
package resolver

import (
	"github.com/InditexTech/xk6-grpcresolver/pkg/grpcresolver"
	"google.golang.org/grpc/resolver"
	"os"
)

// Register the custom gRPC resolver.
func init() {
	if err := grpcresolver.LoadSettings(); err != nil {
		grpcresolver.Logger.Error("failed to load xk6-grpcresolver settings: ", err)
		os.Exit(1)
	}

	resolver.Register(&grpcresolver.Builder{})
}
