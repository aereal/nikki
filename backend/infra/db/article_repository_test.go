package db_test

import (
	"errors"
	"testing"
	"time"

	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/infra/db/test"
	"github.com/google/go-cmp/cmp"
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
	articleID1 := articleRepo.ArticleIDGenerator.GenerateID()
	articleID2 := articleRepo.ArticleIDGenerator.GenerateID()
	aggregate := &domain.ImportArticlesAggregate{
		Articles: []*domain.ArticleToImport{
			{
				ArticleID:         articleID1,
				Slug:              "article_1",
				ArticleRevisionID: articleRepo.ArticleRevisionIDGenerator.GenerateID(),
				Title:             "title 1",
				Body:              "<p>body</p>",
				AuthoredAt:        time.Now(),
				Categories:        categories,
			},
			{
				ArticleID:         articleID2,
				Slug:              "article_2",
				ArticleRevisionID: articleRepo.ArticleRevisionIDGenerator.GenerateID(),
				Title:             "title 2",
				Body:              "<p>body</p>",
				AuthoredAt:        time.Now(),
				Categories:        categories[1:],
			},
		},
	}
	if err := articleRepo.ImportArticles(t.Context(), aggregate); err != nil {
		t.Fatal(err)
	}

	gotArticle, err := articleRepo.FindArticleBySlug(t.Context(), aggregate.Articles[0].Slug)
	if err != nil {
		t.Fatal(err)
	}
	wantArticle := &domain.Article{
		ArticleID: articleID1,
		Slug:      "article_1",
		Title:     "title 1",
	}
	if diff := cmp.Diff(wantArticle, gotArticle); diff != "" {
		t.Errorf("article (-want, +got):\n%s", diff)
	}

	if _, gotErr := articleRepo.FindArticleBySlug(t.Context(), "not_found"); !errors.Is(gotErr, domain.ArticleBySlugNotFound("not_found")) {
		t.Errorf("unexpected error: (%T) %s", gotErr, gotErr)
	}
}
