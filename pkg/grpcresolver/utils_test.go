// SPDX-FileCopyrightText: 2025 INDUSTRIA DE DISEÑO TEXTIL S.A. (INDITEX S.A.)
//
// SPDX-License-Identifier: AGPL-3.0-only

package grpcresolver

import (
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

type Array[T any] struct {
	values []T
}

func (a *Array[T]) Append(value T) {
	a.values = append(a.values, value)
}

// TestClientConnImpl struct that implements resolver.ClientConn
// Copied from google.golang.org/grpc@v1.60.0/internal/testutils/resolver.go
// The methods cannot be implemented with pointer receivers.
type TestClientConnImpl struct {
	stateUpdates *Array[resolver.State]
}

func (t TestClientConnImpl) UpdateState(state resolver.State) error {
	t.stateUpdates.Append(state)
	return nil
}

func (t TestClientConnImpl) ReportError(_ error) {}

func (t TestClientConnImpl) NewAddress(_ []resolver.Address) {}

func (t TestClientConnImpl) ParseServiceConfig(_ string) *serviceconfig.ParseResult {
	return &serviceconfig.ParseResult{}
}
