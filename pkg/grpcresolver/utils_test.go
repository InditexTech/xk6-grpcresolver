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
type TestClientConnImpl struct {
	//stateUpdates atomic.Pointer[[]resolver.State]
	stateUpdates *Array[resolver.State]
}

func (t TestClientConnImpl) UpdateState(state resolver.State) error {
	//var stateUpdates []resolver.State
	//if stateUpdatesPtr := t.stateUpdates.Load(); stateUpdatesPtr != nil {
	//	stateUpdates = *stateUpdatesPtr
	//}
	//
	//stateUpdates = append(stateUpdates, state)
	//t.stateUpdates.Store(&stateUpdates)
	//return nil

	t.stateUpdates.Append(state)
	return nil
}

func (t TestClientConnImpl) ReportError(_ error) {}

func (t TestClientConnImpl) NewAddress(_ []resolver.Address) {}

func (t TestClientConnImpl) ParseServiceConfig(_ string) *serviceconfig.ParseResult {
	return &serviceconfig.ParseResult{}
}
