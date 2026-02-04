//go:build wireinject

package test

import (
	"context"

	"github.com/aereal/nikki/backend/infra/db"
	"github.com/aereal/nikki/backend/infra/db/exec"
	"github.com/aereal/nikki/backend/o11y"
	"github.com/google/wire"
)

func NewTestCategoryRepository(_ context.Context) (*TestCategoryRepository, error) {
	wire.Build(
		db.ProvideCategoryRepository,
		db.ProvideDB,
		db.ProvideMemoryEndpoint,
		exec.ProvideRunner,
		o11y.ProvideNoopTracerProvider,
		provideProvisionedDB,
		provideTestCategoryRepository,
		wire.Bind(new(db.Endpoint), new(*db.MemoryEndpoint)),
		wire.Bind(new(exec.Context), new(*exec.Runner)),
	)
	return nil, nil
}

func NewTestArticleRepository(_ context.Context) (*TestArticleRepository, error) {
	wire.Build(
		db.ProvideArticleRepository,
		db.ProvideDB,
		db.ProvideMemoryEndpoint,
		exec.ProvideRunner,
		o11y.ProvideNoopTracerProvider,
		provideProvisionedDB,
		provideTestArticleRepository,
		wire.Bind(new(db.Endpoint), new(*db.MemoryEndpoint)),
		wire.Bind(new(exec.Context), new(*exec.Runner)),
	)
	return nil, nil
}
