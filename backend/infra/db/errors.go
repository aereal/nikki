package db

var ErrEmptyFile EmptyFileError

type EmptyFileError struct{}

var _ error = EmptyFileError{}

func (EmptyFileError) Error() string { return "empty file path" }

var ErrNoValuesToInsert NoValuesToInsertError

type NoValuesToInsertError struct{}

func (NoValuesToInsertError) Error() string { return "no values to insert" }
