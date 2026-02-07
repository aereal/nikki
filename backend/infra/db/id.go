package db

import (
	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/usecases/ports"
	"github.com/google/wire"
	"github.com/rs/xid"
)

var (
	ArticleIDGeneratorProvider = wire.NewSet(
		ProvideArticleIDGenerator,
		wire.Bind(new(ports.ArticleIDGenerator), new(IDGenerator[domain.ArticleID])),
	)
	ArticleRevisionIDGeneratorProvider = wire.NewSet(
		ProvideArticleRevisionIDGenerator,
		wire.Bind(new(ports.ArticleRevisionIDGenerator), new(IDGenerator[domain.ArticleRevisionID])),
	)
	CategoryIDGeneratorProvider = wire.NewSet(
		ProvideCategoryIDGenerator,
		wire.Bind(new(ports.CategoryIDGenerator), new(IDGenerator[domain.CategoryID])),
	)
)

func ProvideArticleIDGenerator() IDGenerator[domain.ArticleID] {
	return IDGenerator[domain.ArticleID]{}
}

func ProvideArticleRevisionIDGenerator() IDGenerator[domain.ArticleRevisionID] {
	return IDGenerator[domain.ArticleRevisionID]{}
}

func ProvideCategoryIDGenerator() IDGenerator[domain.CategoryID] {
	return IDGenerator[domain.CategoryID]{}
}

type IDGenerator[ID ~string] struct{}

var _ ports.IDGenerator[string] = IDGenerator[string]{}

func (IDGenerator[ID]) GenerateID() ID { return ID(xid.New().String()) }
