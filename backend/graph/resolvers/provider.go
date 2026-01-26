//go:generate bash ./generate.bash

package resolvers

func ProvideResolver() *Resolver {
	return &Resolver{}
}
