package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"iter"
	"time"

	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/infra/db/dto"
	"github.com/aereal/nikki/backend/infra/db/exec"
	"github.com/aereal/nikki/backend/infra/db/queries"
	"github.com/aereal/nikki/backend/o11y"
	"github.com/aereal/nikki/backend/utils"
	"github.com/aereal/optional"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func ProvideArticleRepository(tp trace.TracerProvider, execCtx exec.Context) *ArticleRepository {
	return &ArticleRepository{
		tracer:  tp.Tracer("github.com/aereal/nikki/backend/infra/db.ArticleRepository"),
		execCtx: execCtx,
	}
}

type ArticleRepository struct {
	tracer  trace.Tracer
	execCtx exec.Context
}

var _ domain.ArticleRepository = (*ArticleRepository)(nil)

func (r *ArticleRepository) FindArticles(ctx context.Context, first int, direction domain.OrderDirection, cursor optional.Option[time.Time]) (_ []*domain.Article, _ optional.Option[time.Time], err error) {
	ctx, span := r.tracer.Start(ctx, "FindArticles", trace.WithAttributes(attribute.Int("first", first), attribute.String("direction", direction.String())))
	defer func() { o11y.FinishSpan(span, err) }()

	iterateArticles, err := r.findArticles(ctx, first, direction, cursor)
	if err != nil {
		return nil, optional.None[time.Time](), err
	}
	ret, nextCursor := utils.Paginate(first, cursorOfArticle, iterateArticles)
	return ret, optional.FromPtr(nextCursor), nil
}

func (r *ArticleRepository) findArticles(ctx context.Context, first int, direction domain.OrderDirection, cursor optional.Option[time.Time]) (iter.Seq[*domain.Article], error) {
	limit := int64(first + 1)
	q := queries.New(r.execCtx)
	switch {
	case direction == domain.OrderDirectionDesc && optional.IsSome(cursor):
		val, _ := optional.Unwrap(cursor)
		rows, err := q.FindLatestArticlesAfter(ctx, queries.FindLatestArticlesAfterParams{Limit: limit, After: dto.DateTime(val)})
		if err != nil {
			return nil, err
		}
		return func(yield func(*domain.Article) bool) {
			for _, row := range rows {
				if !yield(row.ToArticle()) {
					return
				}
			}
		}, nil
	case direction == domain.OrderDirectionDesc && optional.IsNone(cursor):
		rows, err := q.FindLatestArticles(ctx, limit)
		if err != nil {
			return nil, err
		}
		return func(yield func(*domain.Article) bool) {
			for _, row := range rows {
				if !yield(row.ToArticle()) {
					return
				}
			}
		}, nil
	case direction == domain.OrderDirectionAsc && optional.IsSome(cursor):
		val, _ := optional.Unwrap(cursor)
		rows, err := q.FindEarlyArticlesBefore(ctx, queries.FindEarlyArticlesBeforeParams{Limit: limit, Before: dto.DateTime(val)})
		if err != nil {
			return nil, err
		}
		return func(yield func(*domain.Article) bool) {
			for _, row := range rows {
				if !yield(row.ToArticle()) {
					return
				}
			}
		}, nil
	case direction == domain.OrderDirectionAsc && optional.IsNone(cursor):
		rows, err := q.FindEarlyArticles(ctx, limit)
		if err != nil {
			return nil, err
		}
		return func(yield func(*domain.Article) bool) {
			for _, row := range rows {
				if !yield(row.ToArticle()) {
					return
				}
			}
		}, nil
	default:
		return nil, fmt.Errorf("unsupported direction=%s after=%v", direction, cursor) //nolint:err113
	}
}

func (r *ArticleRepository) FindArticleBySlug(ctx context.Context, slug string) (_ *domain.Article, err error) {
	ctx, span := r.tracer.Start(ctx, "FindArticleBySlug", trace.WithAttributes(attribute.String("slug", slug)))
	defer func() { o11y.FinishSpan(span, err) }()

	article, err := queries.New(r.execCtx).FindArticleBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ArticleBySlugNotFound(slug)
		}
		return nil, err
	}
	return &domain.Article{
		ArticleID:   article.ArticleID,
		Slug:        article.Slug,
		Title:       article.Title,
		Body:        article.Body,
		PublishedAt: time.Time(article.PublishedAt),
	}, nil
}

func (r *ArticleRepository) ImportArticles(ctx context.Context, aggregate *domain.ImportArticlesAggregate) (err error) {
	ctx, span := r.tracer.Start(ctx, "ImportArticles")
	defer func() { o11y.FinishSpan(span, err) }()

	if len(aggregate.Articles) == 0 {
		return ErrNoValuesToInsert
	}

	if err := r.createArticles(ctx, aggregate.Articles); err != nil {
		return err
	}
	if err := r.createRevisions(ctx, aggregate.Articles); err != nil {
		return err
	}
	if err := r.createArticlePublications(ctx, aggregate.Articles); err != nil {
		return err
	}
	if err := r.createArticleCategoryMappings(ctx, aggregate.Articles); err != nil {
		return err
	}
	return nil
}

func (r *ArticleRepository) createArticles(ctx context.Context, articles []*domain.ArticleToImport) error {
	params := make(queries.BulkCreateArticlesParams, len(articles))
	for i, a := range articles {
		params[i].ArticleID = a.ArticleID
		params[i].Slug = a.Slug
	}
	return queries.New(r.execCtx).BulkCreateArticles(ctx, params)
}

func (r *ArticleRepository) createRevisions(ctx context.Context, articles []*domain.ArticleToImport) error {
	params := make(queries.BulkCreateArticleRevisionsParams, len(articles))
	for i, a := range articles {
		params[i].ArticleID = a.ArticleID
		params[i].ArticleRevisionID = a.ArticleRevisionID
		params[i].Title = a.Title
		params[i].Body = a.Body
		params[i].AuthoredAt = dto.DateTime(a.AuthoredAt)
	}
	return queries.New(r.execCtx).BulkCreateArticleRevisions(ctx, params)
}

func (r *ArticleRepository) createArticlePublications(ctx context.Context, articles []*domain.ArticleToImport) error {
	params := make(queries.BulkCreateArticlePublicationsParams, len(articles))
	for i, a := range articles {
		params[i].ArticleID = a.ArticleID
		params[i].ArticleRevisionID = a.ArticleRevisionID
		params[i].PublishedAt = dto.DateTime(a.AuthoredAt)
	}
	return queries.New(r.execCtx).BulkCreateArticlePublications(ctx, params)
}

func (r *ArticleRepository) createArticleCategoryMappings(ctx context.Context, articles []*domain.ArticleToImport) error {
	params := make(queries.BulkMapArticleCategoryParams, 0)
	for _, a := range articles {
		for _, c := range a.Categories {
			params = append(params, queries.MapArticleCategoryParams{
				ArticleID:  a.ArticleID,
				CategoryID: c.CategoryID,
			})
		}
	}
	if len(params) == 0 {
		return nil
	}
	return queries.New(r.execCtx).BulkMapArticleCategory(ctx, params)
}

func cursorOfArticle(a *domain.Article) time.Time {
	return a.PublishedAt
}
