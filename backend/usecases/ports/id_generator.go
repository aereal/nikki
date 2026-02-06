package ports

import "github.com/aereal/nikki/backend/domain"

type IDGenerator[ID any] interface {
	GenerateID() ID
}

type ArticleIDGenerator = IDGenerator[domain.ArticleID]

type ArticleRevisionIDGenerator = IDGenerator[domain.ArticleRevisionID]
