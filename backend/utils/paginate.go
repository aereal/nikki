package utils

import "iter"

func Paginate[T any, C any](requestSize int, calculateCursor func(T) C, values iter.Seq[T]) ([]T, *C) {
	ret := make([]T, 0, requestSize)
	for v := range values {
		if len(ret) == requestSize {
			cursor := calculateCursor(v)
			return ret, &cursor
		}
		ret = append(ret, v)
	}
	return ret, nil
}
