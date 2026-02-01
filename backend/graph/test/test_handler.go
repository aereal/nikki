package test

import (
	"github.com/aereal/nikki/backend/graph"
)

func provideHandler(h graph.Handler) *Handler {
	return &Handler{
		Handler: h,
	}
}

type Handler struct {
	graph.Handler
}
