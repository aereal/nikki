//go:generate bash ./generate.bash

package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

func ProvideResolver() *Resolver {
	return &Resolver{}
}

type Resolver struct{}
