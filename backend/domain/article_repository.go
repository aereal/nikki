package domain

import (
	"context"
	"time"

	"github.com/aereal/optional"
)

type ArticleRepository interface {
	ImportArticles(ctx context.Context, aggregate *ImportArticlesAggregate) error
	FindArticleBySlug(ctx context.Context, slug string) (*Article, error)
	FindArticles(ctx context.Context, first int, direction OrderDirection, cursor optional.Option[time.Time]) ([]*Article, optional.Option[time.Time], error)
}
