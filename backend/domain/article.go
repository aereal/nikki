package domain

type ArticleID string

type Article struct {
	ArticleID ArticleID
	Slug      string
	Title     string
}
