package env_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aereal/nikki/backend/env"
)

func TestMissingEnvironmentVariableError_Error(t *testing.T) {
	t.Parallel()

	err := &env.MissingEnvironmentVariableError{Name: "k1"}
	if got := err.Error(); got != "missing environment variable: k1" {
		t.Errorf("invalid Error(): %s", got)
	}
}

func TestMissingEnvironmentVariableError_Is(t *testing.T) {
	t.Parallel()

	lhs := &env.MissingEnvironmentVariableError{Name: "k1"}
	testCases := []struct {
		rhs  error
		want bool
	}{
		{rhs: &env.MissingEnvironmentVariableError{Name: "k1"}, want: true},
		{rhs: &env.MissingEnvironmentVariableError{Name: "k2"}, want: true},
		{rhs: errors.New("missing environment variable: k1"), want: false}, //nolint:err113 // allow dynamic test in the test
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%#v", tc.rhs), func(t *testing.T) {
			t.Parallel()

			got := errors.Is(lhs, tc.rhs)
			if got != tc.want {
				t.Errorf("want=%v got=%v", tc.want, got)
			}
		})
	}
}
