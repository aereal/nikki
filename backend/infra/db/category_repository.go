package db

import (
	"context"

	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/infra/db/exec"
	"github.com/aereal/nikki/backend/infra/db/queries"
	"github.com/aereal/nikki/backend/o11y"
	"go.opentelemetry.io/otel/trace"
)

func ProvideCategoryRepository(tp trace.TracerProvider, execCtx exec.Context, idGen IDGenerator[domain.CategoryID]) *CategoryRepository {
	return &CategoryRepository{
		execCtx:     execCtx,
		tracer:      tp.Tracer("github.com/aereal/nikki/backend/infra/db.CategoryRepository"),
		idGenerator: idGen,
	}
}

type CategoryRepository struct {
	tracer      trace.Tracer
	execCtx     exec.Context
	idGenerator IDGenerator[domain.CategoryID]
}

var _ domain.CategoryRepository = (*CategoryRepository)(nil)

func (r *CategoryRepository) ImportCategories(ctx context.Context, names []string) (err error) {
	ctx, span := r.tracer.Start(ctx, "ImportCategories")
	defer func() { o11y.FinishSpan(span, err) }()

	if len(names) == 0 {
		return ErrNoValuesToInsert
	}

	params := queries.BulkImportCategoriesParams{}
	for _, name := range names {
		params = append(params, queries.ImportCategoriesParams{
			CategoryID: r.idGenerator.GenerateID(),
			Name:       name,
		})
	}

	if err := queries.New(r.execCtx).BulkImportCategories(ctx, params); err != nil {
		return err
	}
	return nil
}

func (r *CategoryRepository) FindCategoriesByNames(ctx context.Context, names []string) (_ []*domain.Category, err error) {
	ctx, span := r.tracer.Start(ctx, "FindCategoriesByNames")
	defer func() { o11y.FinishSpan(span, err) }()

	records, err := queries.New(r.execCtx).FindCategoriesByNames(ctx, names)
	if err != nil {
		return nil, err
	}
	ret := make([]*domain.Category, len(records))
	for i, record := range records {
		ret[i] = &domain.Category{
			CategoryID: record.CategoryID,
			Name:       record.Name,
		}
	}
	return ret, nil
}
