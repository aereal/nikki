package interactions_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/testutils"
	"github.com/aereal/nikki/backend/usecases/interactions"
	"github.com/aereal/nikki/backend/usecases/interactions/test"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"
)

func TestImportMTExport_ImportMTExport(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		filename string
		doMock   func(t *testing.T, root *test.TestImportMTExport)
		wantErr  error
	}{
		{
			name:     "ok",
			filename: "./testdata/ok.txt",
			doMock: func(t *testing.T, root *test.TestImportMTExport) {
				t.Helper()

				root.ArticleIDGenerator.EXPECT().GenerateID().Return("article-1").Times(1)
				root.ArticleIDGenerator.EXPECT().GenerateID().Return("article-2").Times(1)
				root.CategoryRepository.EXPECT().
					ImportCategories(gomock.Any(), gomock.InAnyOrder([]string{"News", "Product"})).
					Return(nil).
					Times(1)
				root.CategoryRepository.EXPECT().
					FindCategoriesByNames(gomock.Any(), gomock.InAnyOrder([]string{"News", "Product"})).
					Return([]*domain.Category{
						{CategoryID: "cat-1", Name: "News"},
						{CategoryID: "cat-2", Name: "Product"},
					}, nil).
					Times(1)
				root.ArticleRevisionIDGenerator.EXPECT().GenerateID().Return("revision-1").Times(1)
				root.ArticleRevisionIDGenerator.EXPECT().GenerateID().Return("revision-2").Times(1)
				root.ArticleRepository.EXPECT().
					ImportArticles(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, gotAggregate *domain.ImportArticlesAggregate) error {
						if diff := cmp.Diff(wantOKAggregate, gotAggregate); diff != "" {
							t.Errorf("ImportArticles aggregate (-want, +got):\n%s", diff)
						}
						return nil
					}).
					Times(1)
			},
			wantErr: nil,
		},
		{
			name:     "not found",
			filename: "./testdata/not_found.txt",
			wantErr:  testutils.LiteralError("open ./testdata/not_found.txt: no such file or directory"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			root := test.NewTestImportMTExport(ctrl, interactions.MTExportFileName(tc.filename))
			if tc.doMock != nil {
				tc.doMock(t, root)
			}
			gotErr := root.ImportMTExport.ImportMTExport(t.Context())
			if !errors.Is(tc.wantErr, gotErr) {
				t.Errorf("error:\n\twant: (%T) %s\n\t got: (%T) %s", tc.wantErr, tc.wantErr, gotErr, gotErr)
			}
		})
	}
}

var (
	wantAuthoredTime = time.Date(2007, time.August, 8, 15, 0, 0, 0, time.FixedZone("Asia/Tokyo", int((time.Hour*9).Seconds())))

	wantOKAggregate = &domain.ImportArticlesAggregate{
		Articles: []*domain.ArticleToImport{
			{
				Article: &domain.Article{
					ArticleID: "article-1",
					Slug:      "filename",
				},
				ArticleRevision: &domain.ArticleRevision{
					ArticleRevisionID: "revision-1",
					ArticleID:         "article-1",
					Title:             "A dummy title",
					Body:              "これは本文です。\nここに追記の本文が表示されます。\n",
					AuthoredAt:        wantAuthoredTime,
				},
				Categories: []*domain.Category{
					{CategoryID: "cat-1", Name: "News"},
					{CategoryID: "cat-2", Name: "Product"},
				},
			},
			{
				Article: &domain.Article{
					ArticleID: "article-2",
					Slug:      "filename",
				},
				ArticleRevision: &domain.ArticleRevision{
					ArticleRevisionID: "revision-2",
					ArticleID:         "article-2",
					Title:             "2件目の記事",
					Body:              "これは2番目の記事の本文です。 これは\n複数行から成ります。\n",
					AuthoredAt:        wantAuthoredTime,
				},
				Categories: []*domain.Category{
					{CategoryID: "cat-1", Name: "News"},
					{CategoryID: "cat-2", Name: "Product"},
				},
			},
		},
	}
)
