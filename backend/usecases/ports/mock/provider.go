//go:generate go tool mockgen -typed -package mock -destination ./mock_gen.go github.com/aereal/nikki/backend/usecases/ports IDGenerator

package mock

import (
	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/usecases/ports"
	"github.com/google/wire"
	"go.uber.org/mock/gomock"
)

func NewArticleIDGenerator(ctrl *gomock.Controller) *MockIDGenerator[domain.ArticleID] {
	return NewMockIDGenerator[domain.ArticleID](ctrl)
}

func NewArticleRevisionIDGenerator(ctrl *gomock.Controller) *MockIDGenerator[domain.ArticleRevisionID] {
	return NewMockIDGenerator[domain.ArticleRevisionID](ctrl)
}

var (
	ArticleIDGeneratorProvider = wire.NewSet(
		NewArticleIDGenerator,
		wire.Bind(new(ports.ArticleIDGenerator), new(*MockIDGenerator[domain.ArticleID])),
	)
	ArticleRevisionIDGeneratorProvider = wire.NewSet(
		NewArticleRevisionIDGenerator,
		wire.Bind(new(ports.ArticleRevisionIDGenerator), new(*MockIDGenerator[domain.ArticleRevisionID])),
	)
)
