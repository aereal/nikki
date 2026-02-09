package types

import "fmt"

func Cast[T any](v any) (T, error) {
	t, ok := v.(T)
	if !ok {
		return *new(T), &UnexpectedTypeError[T]{Actual: v}
	}
	return t, nil
}

type UnexpectedTypeError[Expected any] struct {
	Actual any
}

func (e *UnexpectedTypeError[Expected]) Error() string {
	return fmt.Sprintf("expected type %T but got %T", *new(Expected), e.Actual)
}
