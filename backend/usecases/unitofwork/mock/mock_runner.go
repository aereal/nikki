package mock

import (
	"context"

	"github.com/aereal/nikki/backend/usecases/unitofwork"
)

func ProvideMockRunner() unitofwork.Runner { return mockRunner{} }

type mockRunner struct{}

var _ unitofwork.Runner = mockRunner{}

func (mockRunner) StartUnitOfWork(ctx context.Context) (context.Context, unitofwork.Finisher, error) {
	return ctx, unitofwork.NoopFinisher, nil
}
