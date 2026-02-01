package graph

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/aereal/nikki/backend/graph/exec"
	"github.com/aereal/nikki/backend/graph/resolvers"
	"github.com/aereal/otelgqlgen"
	"go.opentelemetry.io/otel/trace"
)

type Handler interface{ http.Handler }

func ProviveHandler(tp trace.TracerProvider, root *resolvers.Resolver) Handler {
	es := exec.NewExecutableSchema(exec.Config{
		Resolvers: root,
	})
	h := handler.New(es)
	h.AddTransport(transport.POST{})
	h.Use(otelgqlgen.New(otelgqlgen.WithTracerProvider(tp), otelgqlgen.ShouldTraceCaptureTimings(false)))
	return h
}
