package domain

import "context"

type ArticleRepository interface {
	ImportArticles(ctx context.Context, aggregate *ImportArticlesAggregate) error
}
