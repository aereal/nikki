package usecases

import "context"

type ImportMTExport interface {
	ImportMTExport(_ context.Context) error
}
