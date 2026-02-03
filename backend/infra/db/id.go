package db

import "github.com/rs/xid"

func GenerateID[T ~string]() T {
	return T(xid.New().String())
}
