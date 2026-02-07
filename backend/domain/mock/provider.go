//go:generate go tool mockgen -typed -package mock -destination ./mock_gen.go github.com/aereal/nikki/backend/domain ArticleRepository,CategoryRepository

package mock

import (
	"github.com/aereal/nikki/backend/domain"
	"github.com/google/wire"
	"go.uber.org/mock/gomock"
)

var (
	ArticleRepositoryProvider = wire.NewSet(
		NewMockArticleRepository,
		wire.Bind(new(domain.ArticleRepository), new(*MockArticleRepository)),
	)
	CategoryRepositoryProvider = wire.NewSet(
		NewMockCategoryRepository,
		wire.Bind(new(domain.CategoryRepository), new(*MockCategoryRepository)),
	)
)

func ProvideArticleRepository(ctrl *gomock.Controller) domain.ArticleRepository {
	return NewMockArticleRepository(ctrl)
}

func ProvideCategoryRepository(ctrl *gomock.Controller) domain.CategoryRepository {
	return NewMockCategoryRepository(ctrl)
}
