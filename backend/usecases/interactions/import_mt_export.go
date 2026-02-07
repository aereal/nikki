package interactions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"slices"

	"github.com/aereal/coll"
	"github.com/aereal/mt"
	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/log/attr"
	"github.com/aereal/nikki/backend/o11y"
	"github.com/aereal/nikki/backend/usecases"
	"github.com/aereal/nikki/backend/usecases/ports"
	"github.com/aereal/nikki/backend/usecases/unitofwork"
	"go.opentelemetry.io/otel/trace"
)

type MTExportFileName string

func ProvideImportMTExport(tp trace.TracerProvider, articleRepo domain.ArticleRepository, categoryRepo domain.CategoryRepository, articleIDGenerator ports.ArticleIDGenerator, articleRevisionIDGenerator ports.ArticleRevisionIDGenerator, runner unitofwork.Runner, fileName MTExportFileName) *ImportMTExport {
	return &ImportMTExport{
		tracer:                     tp.Tracer("github.com/aereal/nikki/backend/usecases/interactions.ImportMTExport"),
		articleRepository:          articleRepo,
		categoryRepository:         categoryRepo,
		articleIDGenerator:         articleIDGenerator,
		articleRevisionIDGenerator: articleRevisionIDGenerator,
		runner:                     runner,
		exportFileName:             fileName,
	}
}

type ImportMTExport struct {
	tracer                     trace.Tracer
	articleRepository          domain.ArticleRepository
	categoryRepository         domain.CategoryRepository
	articleIDGenerator         ports.ArticleIDGenerator
	articleRevisionIDGenerator ports.ArticleRevisionIDGenerator
	runner                     unitofwork.Runner
	exportFileName             MTExportFileName
}

var _ usecases.ImportMTExport = (*ImportMTExport)(nil)

func (i *ImportMTExport) ImportMTExport(ctx context.Context) (err error) {
	ctx, span := i.tracer.Start(ctx, "ImportMTExport")
	defer func() { o11y.FinishSpan(span, err) }()

	index := -1
	entries := map[domain.ArticleID]*mt.Entry{}
	f, err := os.Open(string(i.exportFileName))
	if err != nil {
		return err
	}
	defer f.Close()
	for entry, err := range mt.Parse(f) {
		index++
		if err != nil {
			slog.WarnContext(ctx, "MT entry parse failure", slog.Int("index", index), attr.Error(err))
			continue
		}
		articleID := i.articleIDGenerator.GenerateID()
		entries[articleID] = entry
	}

	ctx, finish, err := i.runner.StartUnitOfWork(ctx)
	if err != nil {
		return fmt.Errorf("unitofwork.Runner.StartUnitOfWork: %w", err)
	}
	defer func() { finish(err) }()

	name2category, err := i.importCategories(ctx, entries)
	if err != nil {
		return err
	}

	errs := make([]error, 0, len(entries))
	aggregate := &domain.ImportArticlesAggregate{}
	for _, articleID := range slices.Sorted(maps.Keys(entries)) {
		entry := entries[articleID]
		articleRevisionID := i.articleRevisionIDGenerator.GenerateID()
		articleToImport, err := ports.ConvertMTEntry(articleID, articleRevisionID, entry, name2category)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		aggregate.Articles = append(aggregate.Articles, articleToImport)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	if err := i.articleRepository.ImportArticles(ctx, aggregate); err != nil {
		return fmt.Errorf("domain.ArticleRepository.ImportArticles: %w", err)
	}
	return nil
}

func (i *ImportMTExport) importCategories(ctx context.Context, entries map[domain.ArticleID]*mt.Entry) (map[string]*domain.Category, error) {
	allNames := coll.NewSet[string]()
	for _, entry := range entries {
		names := ports.CategoryNamesOfMTEntry(entry)
		if names.Len() > 0 {
			allNames = allNames.Union(names)
		}
	}
	if allNames.Len() == 0 {
		return map[string]*domain.Category{}, nil
	}

	categoryNames := slices.Collect(allNames.Values())
	if err := i.categoryRepository.ImportCategories(ctx, categoryNames); err != nil {
		return nil, fmt.Errorf("domain.CategoryRepository.ImportCategories: %w", err)
	}
	cats, err := i.categoryRepository.FindCategoriesByNames(ctx, categoryNames)
	if err != nil {
		return nil, fmt.Errorf("domain.CategoryRepository.FindCategoriesByNames: %w", err)
	}
	ret := map[string]*domain.Category{}
	for _, c := range cats {
		ret[c.Name] = c
	}
	return ret, nil
}
