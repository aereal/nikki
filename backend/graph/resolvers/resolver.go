package resolvers

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"

	"github.com/aereal/nikki/backend/graph/dto"
	"github.com/aereal/nikki/backend/graph/exec"
)

type Resolver struct{}

func (r *queryResolver) Article(ctx context.Context, slug string) (*dto.Article, error) {
	return &dto.Article{Slug: slug}, nil
}

func (r *Resolver) Query() exec.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
