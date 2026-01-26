//go:build wireinject

package test

import (
	"github.com/aereal/nikki/backend/graph"
	"github.com/aereal/nikki/backend/graph/resolvers"
	"github.com/aereal/nikki/backend/o11y"
	"github.com/google/wire"
)

func NewHandler() *Handler {
	wire.Build(
		graph.ProviveHandler,
		o11y.ProvideNoopTracerProvider,
		provideHandler,
		resolvers.ProvideResolver,
	)
	return nil
}
