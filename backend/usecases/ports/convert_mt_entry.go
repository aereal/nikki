package ports

import (
	"slices"

	"github.com/aereal/mt"
	"github.com/aereal/nikki/backend/domain"
)

func ConvertMTEntry(articleID domain.ArticleID, articleRevisionID domain.ArticleRevisionID, entry *mt.Entry, name2category map[string]*domain.Category) (*domain.ArticleToImport, error) {
	if err := validateEntry(entry); err != nil {
		return nil, asConvertMTEntryError(articleID, articleRevisionID, err)
	}

	article := &domain.Article{
		ArticleID: articleID,
		Slug:      entry.Basename,
	}
	revision := &domain.ArticleRevision{
		ArticleID:         articleID,
		ArticleRevisionID: articleRevisionID,
		Title:             entry.Title,
		Body:              entry.Body + entry.ExtendedBody,
		AuthoredAt:        entry.Date,
	}
	articleToImport := &domain.ArticleToImport{Article: article, ArticleRevision: revision}
	categoryNames := CategoryNamesOfMTEntry(entry)
	errs := make([]error, 0, categoryNames.Len())
	for _, name := range slices.Sorted(categoryNames.Values()) {
		cat, ok := name2category[name]
		if !ok {
			errs = append(errs, domain.CategoryByNameNotFound(name))
			continue
		}
		articleToImport.Categories = append(articleToImport.Categories, cat)
	}
	if len(errs) > 0 {
		return nil, asConvertMTEntryError(articleID, articleRevisionID, errs...)
	}
	return articleToImport, nil
}

func validateEntry(entry *mt.Entry) *InvalidMTExportEntryError {
	errs := make([]error, 0)
	switch entry.ConvertBreaks {
	case mt.ConvertBreaksNone, mt.ConvertBreaksRichtext: // supported
	default:
		errs = append(errs, &UnsupportedConvertBreaksError{Value: entry.ConvertBreaks})
	}
	if entry.Basename == "" {
		errs = append(errs, ErrEmptyBasename)
	}
	if entry.Date.IsZero() {
		errs = append(errs, ErrEmptyDate)
	}

	if len(errs) > 0 {
		return &InvalidMTExportEntryError{errs: errs}
	}
	return nil
}
