package domain

import "context"

type ArticleRepository interface {
	ImportArticles(ctx context.Context, aggregate *ImportArticlesAggregate) error
	FindArticleBySlug(ctx context.Context, slug string) (*Article, error)
}
