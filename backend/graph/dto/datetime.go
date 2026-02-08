package dto

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/aereal/nikki/backend/types"
)

var parser = types.DateTimeParser("2006-01-02T15:04:05.999-07:00")

func MarshalDateTime(t time.Time) graphql.ContextMarshaler {
	return graphql.ContextWriterFunc(func(_ context.Context, w io.Writer) error {
		return json.NewEncoder(w).Encode(parser.Format(t))
	})
}

func UnmarshalDateTime(_ context.Context, v any) (time.Time, error) {
	sv, err := types.Cast[string](v)
	if err != nil {
		return time.Time{}, err
	}
	return parser.Parse(sv)
}
