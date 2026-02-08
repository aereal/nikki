package resolvers

import (
	"context"

	"github.com/aereal/nikki/backend/graph/dto"
	"github.com/aereal/nikki/backend/graph/exec"
)

func (r *queryResolver) Article(ctx context.Context, slug string) (*dto.Article, error) {
	return &dto.Article{Slug: slug}, nil
}

func (r *Resolver) Query() exec.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
