package utils_test

import (
	"slices"
	"strconv"
	"testing"

	"github.com/aereal/nikki/backend/utils"
	"github.com/google/go-cmp/cmp"
)

func TestPaginate(t *testing.T) {
	t.Parallel()

	cursor := "4"
	testCases := []struct {
		name       string
		values     []int
		size       int
		wantValues []int
		wantCursor *string
	}{
		{
			name:       "empty",
			values:     []int{},
			size:       3,
			wantValues: []int{},
			wantCursor: nil,
		},
		{
			name:       "request > size",
			values:     []int{1, 2},
			size:       3,
			wantValues: []int{1, 2},
			wantCursor: nil,
		},
		{
			name:       "request == size",
			values:     []int{1, 2, 3},
			size:       3,
			wantValues: []int{1, 2, 3},
			wantCursor: nil,
		},
		{
			name:       "request < size",
			values:     []int{1, 2, 3, 4},
			size:       3,
			wantValues: []int{1, 2, 3},
			wantCursor: &cursor,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotValues, gotCursor := utils.Paginate(tc.size, strconv.Itoa, slices.Values(tc.values))
			if diff := cmp.Diff(tc.wantValues, gotValues); diff != "" {
				t.Errorf("values (-want, +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantCursor, gotCursor); diff != "" {
				t.Errorf("cursor (-want, +got):\n%s", diff)
			}
		})
	}
}
