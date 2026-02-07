package test

import "github.com/aereal/nikki/backend/infra/db"

func provideTestCategoryRepository(_ *provisionedDB, r *db.CategoryRepository) *TestCategoryRepository {
	return &TestCategoryRepository{CategoryRepository: r}
}

type TestCategoryRepository struct {
	*db.CategoryRepository
}
