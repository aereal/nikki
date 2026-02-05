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

	categoryRepo, err := test.NewTestCategoryRepository(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	categoryNames := []string{"a", "b"}
	if err := categoryRepo.ImportCategories(t.Context(), categoryNames); err != nil {
		t.Fatal(err)
	}
	categories, err := categoryRepo.FindCategoriesByNames(t.Context(), categoryNames)
	if err != nil {
		t.Fatal(err)
	}

	articleRepo, err := test.NewTestArticleRepository(t.Context())
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
				Categories: categories,
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
				Categories: categories[1:],
			},
		},
	}
	if err := articleRepo.ImportArticles(t.Context(), aggregate); err != nil {
		t.Fatal(err)
	}
}
