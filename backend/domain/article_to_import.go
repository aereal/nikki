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
	Status            ArticleStatus
}

type ArticleStatus int

const (
	ArticleStatusInvalid ArticleStatus = iota
	ArticleStatusDraft
	ArticleStatusPublic
)
