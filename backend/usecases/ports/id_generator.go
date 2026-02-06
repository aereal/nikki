package ports

import "github.com/aereal/nikki/backend/domain"

type IDish interface {
	~string
}

type IDGenerator[ID IDish] interface {
	GenerateID() ID
}

type ArticleIDGenerator = IDGenerator[domain.ArticleID]

type ArticleRevisionIDGenerator = IDGenerator[domain.ArticleRevisionID]

type CategoryIDGenerator = IDGenerator[domain.CategoryID]
