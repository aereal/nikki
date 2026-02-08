package resolvers

import (
	"context"

	"github.com/aereal/nikki/backend/graph/dto"
	"github.com/aereal/nikki/backend/graph/exec"
)

func (r *queryResolver) Article(ctx context.Context, slug string) (*dto.Article, error) {
	article, err := r.articleRepository.FindArticleBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return &dto.Article{
		Slug:  article.Slug,
		Title: article.Title,
	}, nil
}

func (r *Resolver) Query() exec.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
