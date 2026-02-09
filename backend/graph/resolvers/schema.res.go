package resolvers

import (
	"context"
	"time"

	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/graph/dto"
	"github.com/aereal/nikki/backend/graph/exec"
	"github.com/aereal/optional"
)

func (r *queryResolver) Article(ctx context.Context, slug string) (*dto.Article, error) {
	article, err := r.articleRepository.FindArticleBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return &dto.Article{
		Slug:        article.Slug,
		Title:       article.Title,
		Body:        article.Body,
		PublishedAt: article.PublishedAt,
	}, nil
}

func (r *queryResolver) Articles(ctx context.Context, first int, order dto.ArticleOrder) (*dto.ArticleConnection, error) {
	var direction domain.OrderDirection
	switch order.Direction {
	case dto.OrderDirectionAsc:
		direction = domain.OrderDirectionAsc
	case dto.OrderDirectionDesc:
		direction = domain.OrderDirectionDesc
	}
	articles, _, err := r.articleRepository.FindArticles(ctx, first, direction, optional.None[time.Time]())
	if err != nil {
		return nil, err
	}
	conn := &dto.ArticleConnection{Nodes: make([]*dto.Article, len(articles))}
	for i, a := range articles {
		conn.Nodes[i] = &dto.Article{
			Slug:        a.Slug,
			Title:       a.Title,
			Body:        a.Body,
			PublishedAt: a.PublishedAt,
		}
	}
	return conn, nil
}

func (r *Resolver) Query() exec.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
