package db_test

import (
	"errors"
	"testing"
	"time"

	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/infra/db/dto"
	"github.com/aereal/nikki/backend/infra/db/test"
	"github.com/aereal/optional"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestArticleRepository_revise(t *testing.T) {
	t.Parallel()

	articleRepo, err := test.NewTestArticleRepository(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	articleID := articleRepo.ArticleIDGenerator.GenerateID()
	slug := "article_1"
	authoredAt := time.Now()
	aggregate := &domain.ImportArticlesAggregate{
		Articles: []*domain.ArticleToImport{
			{
				ArticleID:         articleID,
				ArticleRevisionID: articleRepo.ArticleRevisionIDGenerator.GenerateID(),
				Slug:              slug,
				Title:             "title 1",
				Body:              "<p>body 1</p>",
				AuthoredAt:        authoredAt,
				Status:            domain.ArticleStatusPublic,
			},
		},
	}
	if err := articleRepo.ImportArticles(t.Context(), aggregate); err != nil {
		t.Fatal(err)
	}

	revisedAt := authoredAt.Add(time.Minute)
	nextRevID := articleRepo.ArticleRevisionIDGenerator.GenerateID()
	param := test.ReviseArticleParam{
		ArticleID:         articleID,
		ArticleRevisionID: nextRevID,
		Title:             "revised title",
		Body:              `<p>revised body</p>`,
		AuthoredAt:        dto.DateTime(revisedAt),
	}
	if err := articleRepo.Revise(t.Context(), []test.ReviseArticleParam{param}); err != nil {
		t.Fatal(err)
	}

	want := &domain.Article{
		ArticleID:   articleID,
		Slug:        slug,
		Title:       "revised title",
		Body:        `<p>revised body</p>`,
		PublishedAt: authoredAt.Truncate(time.Millisecond),
	}
	{
		got, err := articleRepo.FindArticleBySlug(t.Context(), slug)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want, +got):\n%s", diff)
		}
	}
	{
		got, _, err := articleRepo.FindArticles(t.Context(), 1, domain.OrderDirectionAsc, optional.None[time.Time]())
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff([]*domain.Article{want}, got); diff != "" {
			t.Errorf("(-want, +got):\n%s", diff)
		}
	}
}

func TestArticleRepository(t *testing.T) { //nolint:tparallel
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
	articleID3 := articleRepo.ArticleIDGenerator.GenerateID()
	authoredAt1 := time.Now()
	authoredAt2 := authoredAt1.Add(time.Second * 5 * -1)
	authoredAt3 := authoredAt2.Add(time.Second * 5 * -1)
	aggregate := &domain.ImportArticlesAggregate{
		Articles: []*domain.ArticleToImport{
			{
				ArticleID:         articleID1,
				Slug:              "article_1",
				ArticleRevisionID: articleRepo.ArticleRevisionIDGenerator.GenerateID(),
				Title:             "title 1",
				Body:              "<p>body 1</p>",
				AuthoredAt:        authoredAt1,
				Categories:        categories,
				Status:            domain.ArticleStatusPublic,
			},
			{
				ArticleID:         articleID2,
				Slug:              "article_2",
				ArticleRevisionID: articleRepo.ArticleRevisionIDGenerator.GenerateID(),
				Title:             "title 2",
				Body:              "<p>body 2</p>",
				AuthoredAt:        authoredAt2,
				Categories:        categories[1:],
				Status:            domain.ArticleStatusPublic,
			},
			{
				ArticleID:         articleID3,
				Slug:              "article_3",
				ArticleRevisionID: articleRepo.ArticleRevisionIDGenerator.GenerateID(),
				Title:             "title 3",
				Body:              "<p>draft</p>",
				AuthoredAt:        authoredAt3,
				Categories:        categories[1:],
				Status:            domain.ArticleStatusDraft,
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
		ArticleID:   articleID1,
		Slug:        "article_1",
		Title:       "title 1",
		Body:        "<p>body 1</p>",
		PublishedAt: authoredAt1.Truncate(time.Millisecond),
	}
	if diff := cmp.Diff(wantArticle, gotArticle); diff != "" {
		t.Errorf("article (-want, +got):\n%s", diff)
	}

	if _, gotErr := articleRepo.FindArticleBySlug(t.Context(), "not_found"); !errors.Is(gotErr, domain.ArticleBySlugNotFound("not_found")) {
		t.Errorf("unexpected error: (%T) %s", gotErr, gotErr)
	}

	t.Run("first=1 order=desc cursor=none", func(t *testing.T) { //nolint:paralleltest
		got, cursor, err := articleRepo.FindArticles(t.Context(), 1, domain.OrderDirectionDesc, optional.None[time.Time]())
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(optional.Some(authoredAt2.Truncate(time.Millisecond)), cursor, cmp.Transformer("Optional", func(o optional.Option[time.Time]) *time.Time { return o.Ptr() }), cmpopts.EquateComparable(time.Time{})); diff != "" {
			t.Errorf("cursor (-want, +got):\n%s", diff)
		}
		want := []*domain.Article{
			{
				ArticleID:   articleID1,
				Slug:        "article_1",
				Title:       "title 1",
				Body:        "<p>body 1</p>",
				PublishedAt: authoredAt1.Truncate(time.Millisecond),
			},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("articles (-want, +got):\n%s", diff)
		}
	})

	t.Run("first=1 order=desc cursor=1", func(t *testing.T) { //nolint:paralleltest
		got, cursor, err := articleRepo.FindArticles(t.Context(), 1, domain.OrderDirectionDesc, optional.Some(authoredAt2))
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(optional.None[time.Time](), cursor, cmp.Transformer("Optional", func(o optional.Option[time.Time]) *time.Time { return o.Ptr() }), cmpopts.EquateComparable(time.Time{})); diff != "" {
			t.Errorf("cursor (-want, +got):\n%s", diff)
		}
		want := []*domain.Article{
			{
				ArticleID:   articleID1,
				Slug:        "article_1",
				Title:       "title 1",
				Body:        "<p>body 1</p>",
				PublishedAt: authoredAt1.Truncate(time.Millisecond),
			},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("articles (-want, +got):\n%s", diff)
		}
	})
}
