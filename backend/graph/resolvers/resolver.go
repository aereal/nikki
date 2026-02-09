//go:generate bash ./generate.bash

package resolvers

import "github.com/aereal/nikki/backend/domain"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

func ProvideResolver(articleRepo domain.ArticleRepository) *Resolver {
	return &Resolver{
		articleRepository: articleRepo,
	}
}

type Resolver struct {
	articleRepository domain.ArticleRepository
}
