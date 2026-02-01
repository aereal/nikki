package db_test

import (
	"errors"
	"testing"

	"github.com/aereal/nikki/backend/infra/db"
)

func TestFileEndpoint_DataSourceName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		endpoint  *db.FileEndpoint
		wantError error
		want      string
	}{
		{
			name:     "no parameters",
			endpoint: &db.FileEndpoint{Path: "test.db"},
			want:     "file://test.db",
		},
		{
			name:     "with parameters",
			endpoint: &db.FileEndpoint{Path: "test.db", Params: &db.ParameterSet{Cache: db.CacheModePrivate}},
			want:     "file://test.db?cache=private",
		},
		{
			name:      "empty path",
			endpoint:  &db.FileEndpoint{Path: ""},
			wantError: db.ErrEmptyFile,
			want:      "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := tc.endpoint.DataSourceName()
			if !errors.Is(tc.wantError, gotErr) {
				t.Errorf("error:\n\twant: %T %s\n\t got: %T %s", tc.wantError, tc.wantError, gotErr, gotErr)
			}
			if got != tc.want {
				t.Errorf("want=%s got=%s", tc.want, got)
			}
		})
	}
}
