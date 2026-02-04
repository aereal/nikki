package test

import "github.com/aereal/nikki/backend/infra/db"

func provideTestArticleRepository(_ *provisionedDB, r *db.ArticleRepository) *TestArticleRepository {
	return &TestArticleRepository{ArticleRepository: r}
}

type TestArticleRepository struct {
	*db.ArticleRepository
}
