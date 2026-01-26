package db_test

import (
	"testing"

	"github.com/aereal/nikki/backend/infra/db"
)

func TestFileEndpoint_DataSourceName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		endpoint *db.FileEndpoint
		want     string
	}{
		{
			name:     "no parameters",
			endpoint: &db.FileEndpoint{Path: "test.db"},
			want:     "file:test.db",
		},
		{
			name:     "with parameters",
			endpoint: &db.FileEndpoint{Path: "test.db", Params: &db.ParameterSet{Cache: db.CacheModePrivate}},
			want:     "file:test.db?cache=private",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.endpoint.DataSourceName()
			if got != tc.want {
				t.Errorf("want=%s got=%s", tc.want, got)
			}
		})
	}
}
