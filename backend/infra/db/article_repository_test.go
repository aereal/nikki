package db_test

import (
	"testing"
	"time"

	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/infra/db"
	"github.com/aereal/nikki/backend/infra/db/test"
)

func TestArticleRepository(t *testing.T) {
	t.Parallel()

	r, err := test.NewTestArticleRepository(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	articleID1 := db.GenerateID[domain.ArticleID]()
	articleID2 := db.GenerateID[domain.ArticleID]()
	aggregate := &domain.ImportArticlesAggregate{
		Articles: []*domain.ArticleToImport{
			{
				Article: &domain.Article{
					ArticleID: articleID1,
					Slug:      "article_1",
				},
				ArticleRevision: &domain.ArticleRevision{
					ArticleID:         articleID1,
					ArticleRevisionID: db.GenerateID[domain.ArticleRevisionID](),
					Title:             "title",
					Body:              "<p>body</p>",
					AuthoredAt:        time.Now(),
				},
			},
			{
				Article: &domain.Article{
					ArticleID: articleID2,
					Slug:      "article_2",
				},
				ArticleRevision: &domain.ArticleRevision{
					ArticleID:         articleID2,
					ArticleRevisionID: db.GenerateID[domain.ArticleRevisionID](),
					Title:             "title",
					Body:              "<p>body</p>",
					AuthoredAt:        time.Now(),
				},
			},
		},
	}
	if err := r.ImportArticles(t.Context(), aggregate); err != nil {
		t.Fatal(err)
	}
}
