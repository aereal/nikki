package unitofwork

import "context"

type Runner interface {
	StartUnitOfWork(ctx context.Context) (context.Context, Finisher, error)
}

type Finisher func(error)

var NoopFinisher = Finisher(func(error) {})
