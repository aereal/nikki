package db

var ErrEmptyFile EmptyFileError

type EmptyFileError struct{}

var _ error = EmptyFileError{}

func (EmptyFileError) Error() string { return "empty file path" }
