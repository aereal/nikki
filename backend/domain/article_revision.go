package domain

import "time"

type ArticleRevisionID string

type ArticleRevision struct {
	ArticleRevisionID ArticleRevisionID
	ArticleID         ArticleID
	Title             string
	Body              string
	AuthoredAt        time.Time
}
