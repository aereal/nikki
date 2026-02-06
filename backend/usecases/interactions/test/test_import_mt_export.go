package test

import (
	"github.com/aereal/nikki/backend/domain"
	domainmock "github.com/aereal/nikki/backend/domain/mock"
	"github.com/aereal/nikki/backend/usecases/interactions"
	portsmock "github.com/aereal/nikki/backend/usecases/ports/mock"
)

func provideTestImportMTExport(i *interactions.ImportMTExport, articleRepo *domainmock.MockArticleRepository, categoryRepo *domainmock.MockCategoryRepository, articleIDGenerator *portsmock.MockIDGenerator[domain.ArticleID], articleRevisionIDGenerator *portsmock.MockIDGenerator[domain.ArticleRevisionID]) *TestImportMTExport {
	return &TestImportMTExport{
		ImportMTExport:             i,
		ArticleRepository:          articleRepo,
		CategoryRepository:         categoryRepo,
		ArticleIDGenerator:         articleIDGenerator,
		ArticleRevisionIDGenerator: articleRevisionIDGenerator,
	}
}

type TestImportMTExport struct {
	*interactions.ImportMTExport

	ArticleRepository          *domainmock.MockArticleRepository
	CategoryRepository         *domainmock.MockCategoryRepository
	ArticleIDGenerator         *portsmock.MockIDGenerator[domain.ArticleID]
	ArticleRevisionIDGenerator *portsmock.MockIDGenerator[domain.ArticleRevisionID]
}
