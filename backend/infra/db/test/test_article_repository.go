package test

import (
	"context"
	"database/sql"

	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/infra/db"
	"github.com/aereal/nikki/backend/infra/db/test/queries"
)

func provideTestArticleRepository(rawdb *provisionedDB, r *db.ArticleRepository, articleIDGen db.IDGenerator[domain.ArticleID], articleRevisionIDGen db.IDGenerator[domain.ArticleRevisionID]) *TestArticleRepository {
	return &TestArticleRepository{
		rawdb:                      rawdb.db,
		ArticleRepository:          r,
		ArticleIDGenerator:         articleIDGen,
		ArticleRevisionIDGenerator: articleRevisionIDGen,
	}
}

type TestArticleRepository struct {
	*db.ArticleRepository

	rawdb                      *sql.DB
	ArticleIDGenerator         db.IDGenerator[domain.ArticleID]
	ArticleRevisionIDGenerator db.IDGenerator[domain.ArticleRevisionID]
}

type ReviseArticleParam = queries.ReviseArticleParams

func (r *TestArticleRepository) Revise(ctx context.Context, params []ReviseArticleParam) error {
	return queries.New(r.rawdb).BulkReviseArticle(ctx, params)
}
