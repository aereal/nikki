package test

import (
	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/infra/db"
)

func provideTestArticleRepository(_ *provisionedDB, r *db.ArticleRepository, articleIDGen db.IDGenerator[domain.ArticleID], articleRevisionIDGen db.IDGenerator[domain.ArticleRevisionID]) *TestArticleRepository {
	return &TestArticleRepository{
		ArticleRepository:          r,
		ArticleIDGenerator:         articleIDGen,
		ArticleRevisionIDGenerator: articleRevisionIDGen,
	}
}

type TestArticleRepository struct {
	*db.ArticleRepository

	ArticleIDGenerator         db.IDGenerator[domain.ArticleID]
	ArticleRevisionIDGenerator db.IDGenerator[domain.ArticleRevisionID]
}
