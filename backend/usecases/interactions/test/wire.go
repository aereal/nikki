//go:build wireinject

package test

import (
	domainmock "github.com/aereal/nikki/backend/domain/mock"
	"github.com/aereal/nikki/backend/o11y"
	"github.com/aereal/nikki/backend/usecases/interactions"
	portsmock "github.com/aereal/nikki/backend/usecases/ports/mock"
	uowmock "github.com/aereal/nikki/backend/usecases/unitofwork/mock"
	"github.com/google/wire"
	"go.uber.org/mock/gomock"
)

func NewTestImportMTExport(_ *gomock.Controller, _ interactions.MTExportFileName) *TestImportMTExport {
	wire.Build(
		domainmock.ArticleRepositoryProvider,
		domainmock.CategoryRepositoryProvider,
		interactions.ProvideImportMTExport,
		o11y.ProvideNoopTracerProvider,
		portsmock.ArticleIDGeneratorProvider,
		portsmock.ArticleRevisionIDGeneratorProvider,
		provideTestImportMTExport,
		uowmock.ProvideMockRunner,
	)
	return nil
}
