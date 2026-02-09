//go:build wireinject

package test

import (
	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/domain/mock"
	"github.com/aereal/nikki/backend/graph"
	"github.com/aereal/nikki/backend/graph/resolvers"
	"github.com/aereal/nikki/backend/o11y"
	"github.com/google/wire"
	"go.uber.org/mock/gomock"
)

func NewHandler(_ *gomock.Controller) *Handler {
	wire.Build(
		graph.ProviveHandler,
		mock.NewMockArticleRepository,
		o11y.ProvideNoopTracerProvider,
		provideHandler,
		resolvers.ProvideResolver,
		wire.Bind(new(domain.ArticleRepository), new(*mock.MockArticleRepository)),
	)
	return nil
}
