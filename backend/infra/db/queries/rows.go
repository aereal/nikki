package queries

import (
	"time"

	"github.com/aereal/nikki/backend/domain"
)

func (row *FindLatestArticlesRow) ToArticle() *domain.Article {
	return &domain.Article{
		ArticleID:   row.ArticleID,
		Slug:        row.Slug,
		Title:       row.Title,
		Body:        row.Body,
		PublishedAt: time.Time(row.PublishedAt),
	}
}

func (row *FindLatestArticlesAfterRow) ToArticle() *domain.Article {
	return &domain.Article{
		ArticleID:   row.ArticleID,
		Slug:        row.Slug,
		Title:       row.Title,
		Body:        row.Body,
		PublishedAt: time.Time(row.PublishedAt),
	}
}

func (row *FindEarlyArticlesRow) ToArticle() *domain.Article {
	return &domain.Article{
		ArticleID:   row.ArticleID,
		Slug:        row.Slug,
		Title:       row.Title,
		Body:        row.Body,
		PublishedAt: time.Time(row.PublishedAt),
	}
}

func (row *FindEarlyArticlesBeforeRow) ToArticle() *domain.Article {
	return &domain.Article{
		ArticleID:   row.ArticleID,
		Slug:        row.Slug,
		Title:       row.Title,
		Body:        row.Body,
		PublishedAt: time.Time(row.PublishedAt),
	}
}
