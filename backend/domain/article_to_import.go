package domain

import "time"

type ArticleToImport struct {
	ArticleID         ArticleID
	Slug              string
	ArticleRevisionID ArticleRevisionID
	Title             string
	Body              string
	AuthoredAt        time.Time
	Categories        []*Category
}
