package ports

import (
	"slices"

	"github.com/aereal/mt"
	"github.com/aereal/nikki/backend/domain"
)

func ConvertMTEntry(articleID domain.ArticleID, articleRevisionID domain.ArticleRevisionID, entry *mt.Entry, name2category map[string]*domain.Category) (*domain.ArticleToImport, error) {
	errs := make([]error, 0)
	if validateErrs := validateEntry(entry); validateErrs != nil {
		errs = append(errs, validateErrs...)
	}


	cats := make([]*domain.Category, 0)
	categoryNames := CategoryNamesOfMTEntry(entry)
	for _, name := range slices.Sorted(categoryNames.Values()) {
		cat, ok := name2category[name]
		if !ok {
			errs = append(errs, domain.CategoryByNameNotFound(name))
			continue
		}
		cats = append(cats, cat)
	}
	if len(errs) > 0 {
		return nil, &ConvertMTEntryError{
			ArticleID:         articleID,
			ArticleRevisionID: articleRevisionID,
			Errs:              errs,
		}
	}
	return &domain.ArticleToImport{
		ArticleID:         articleID,
		Slug:              entry.Basename,
		ArticleRevisionID: articleRevisionID,
		Title:             entry.Title,
		Body:              entry.Body + entry.ExtendedBody,
		AuthoredAt:        entry.Date,
		Categories:        cats,
	}, nil
}

func validateEntry(entry *mt.Entry) []error {
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
	return errs
}
