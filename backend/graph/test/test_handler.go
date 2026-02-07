package test

import (
	"github.com/aereal/nikki/backend/domain/mock"
	"github.com/aereal/nikki/backend/graph"
)

func provideHandler(h graph.Handler, articleRepo *mock.MockArticleRepository) *Handler {
	return &Handler{
		Handler:           h,
		ArticleRepository: articleRepo,
	}
}

type Handler struct {
	graph.Handler

	ArticleRepository *mock.MockArticleRepository
}
