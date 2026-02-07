package ports_test

import (
	"errors"
	"testing"
	"time"

	"github.com/aereal/mt"
	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/usecases/ports"
	"github.com/google/go-cmp/cmp"
)

func TestConvertMTEntry(t *testing.T) {
	t.Parallel()

	date := time.Date(2018, time.February, 3, 12, 34, 56, 0, time.UTC)

	testCases := []struct {
		name      string
		entry     *mt.Entry
		mapping   map[string]*domain.Category
		wantValue *domain.ArticleToImport
		wantErr   error
	}{
		{
			name: "ok",
			entry: &mt.Entry{
				ConvertBreaks:   mt.ConvertBreaksNone,
				Basename:        "basename",
				Date:            date,
				PrimaryCategory: "a",
				Category:        []string{"b", "c"},
				Title:           "title",
				Body:            "<p>body</p>",
				ExtendedBody:    "<p>extended</p>",
			},
			mapping: map[string]*domain.Category{
				"a": {CategoryID: "1", Name: "a"},
				"b": {CategoryID: "2", Name: "b"},
				"c": {CategoryID: "3", Name: "c"},
			},
			wantValue: &domain.ArticleToImport{
				Article: &domain.Article{
					ArticleID: "100",
					Slug:      "basename",
				},
				ArticleRevision: &domain.ArticleRevision{
					ArticleID:         "100",
					ArticleRevisionID: "1",
					Title:             "title",
					Body:              "<p>body</p><p>extended</p>",
					AuthoredAt:        date,
				},
				Categories: []*domain.Category{
					{CategoryID: "1", Name: "a"},
					{CategoryID: "2", Name: "b"},
					{CategoryID: "3", Name: "c"},
				},
			},
			wantErr: nil,
		},
		{
			name: "category not given",
			entry: &mt.Entry{
				ConvertBreaks:   mt.ConvertBreaksNone,
				Basename:        "basename",
				Date:            time.Now(),
				PrimaryCategory: "cat-a",
			},
			mapping:   map[string]*domain.Category{},
			wantValue: nil,
			wantErr:   domain.CategoryByNameNotFound("cat-a"),
		},
		{
			name: "unsupported convert breaks",
			entry: &mt.Entry{
				ConvertBreaks: mt.ConvertBreaksMarkdownWithSmartyPants,
				Basename:      "basename",
				Date:          time.Now(),
			},
			mapping:   map[string]*domain.Category{},
			wantValue: nil,
			wantErr:   ports.WrapInvalidMTExportEntryError(&ports.UnsupportedConvertBreaksError{Value: mt.ConvertBreaksMarkdownWithSmartyPants}),
		},
		{
			name: "empty basename",
			entry: &mt.Entry{
				ConvertBreaks: mt.ConvertBreaksNone,
				Date:          time.Now(),
			},
			mapping:   map[string]*domain.Category{},
			wantValue: nil,
			wantErr:   ports.WrapInvalidMTExportEntryError(ports.ErrEmptyBasename),
		},
		{
			name: "empty date",
			entry: &mt.Entry{
				Basename:      "basename",
				ConvertBreaks: mt.ConvertBreaksNone,
			},
			mapping:   map[string]*domain.Category{},
			wantValue: nil,
			wantErr:   ports.WrapInvalidMTExportEntryError(ports.ErrEmptyDate),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := ports.ConvertMTEntry("100", "1", tc.entry, tc.mapping)
			if !errors.Is(tc.wantErr, gotErr) {
				t.Errorf("error:\n\twant: (%T) %s\n\t got: (%T) %s", tc.wantErr, tc.wantErr, gotErr, gotErr)
			}
			if diff := cmp.Diff(tc.wantValue, got); diff != "" {
				t.Errorf("(-want, +got):\n%s", diff)
			}
		})
	}
}
