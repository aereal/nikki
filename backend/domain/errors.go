package domain

import (
	"errors"
	"fmt"
)

func CategoryByNameNotFound(categoryName string) *NotFoundError[*Category, string] {
	return &NotFoundError[*Category, string]{Key: categoryName}
}

type NotFoundError[Value any, Key comparable] struct {
	Key Key
}

var _ error = (*NotFoundError[any, any])(nil)

func (e *NotFoundError[Value, Key]) Error() string {
	return fmt.Sprintf("%T not found", *new(Value))
}

func (e *NotFoundError[Value, Key]) Is(other error) bool {
	otherNotFoundErr := new(NotFoundError[Value, Key])
	return errors.As(other, &otherNotFoundErr)
}
