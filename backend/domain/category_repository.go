package domain

import "context"

type CategoryRepository interface {
	ImportCategories(ctx context.Context, names []string) error
	FindCategoriesByNames(ctx context.Context, names []string) ([]*Category, error)
}
