package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/infra/db/dto"
	"github.com/aereal/nikki/backend/infra/db/exec"
	"github.com/aereal/nikki/backend/infra/db/queries"
	"github.com/aereal/nikki/backend/o11y"
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
