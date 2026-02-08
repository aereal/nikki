package domain

import "time"

type ArticleID string

type Article struct {
	ArticleID   ArticleID
	Slug        string
	Title       string
	Body        string
	PublishedAt time.Time
}
