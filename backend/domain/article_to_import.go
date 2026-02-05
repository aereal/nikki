package domain

type ArticleToImport struct {
	*Article
	*ArticleRevision

	Categories []*Category
}
