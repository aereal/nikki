package db_test

import (
	"testing"

	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/infra/db/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCategoryRepository_import_find(t *testing.T) {
	t.Parallel()

	names := []string{"a", "b", "c"}
	r, err := test.NewTestCategoryRepository(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	{
		got, err := r.FindCategoriesByNames(t.Context(), names)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff([]*domain.Category{}, got); diff != "" {
			t.Errorf("records (-want, +got):\n%s", diff)
		}
	}

	if err := r.ImportCategories(t.Context(), []string{"a", "b", "c"}); err != nil {
		t.Fatal(err)
	}

	want := []*domain.Category{
		{Name: "a"},
		{Name: "b"},
		{Name: "c"},
	}
	{
		got, err := r.FindCategoriesByNames(t.Context(), names)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(domain.Category{}, "CategoryID")); diff != "" {
			t.Errorf("records (-want, +got):\n%s", diff)
		}
	}
}
