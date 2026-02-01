package env

import (
	"os"
	"strings"
)

func ProvideVariables() Variables {
	vs := Variables{}
	for _, pair := range os.Environ() {
		k, v, ok := strings.Cut(pair, "=")
		if !ok {
			continue
		}
		vs[k] = v
	}
	return vs
}

type Variables map[string]string

type parser[T any] func(string) (T, error)

type scanner[T any] func(vars Variables, name string, ptr *T) error

func scanString(vars Variables, name string, ptr *string) error {
	if v, ok := vars[name]; ok {
		*ptr = v
		return nil
	}
	return &MissingEnvironmentVariableError{Name: name}
}

var _ scanner[string] = scanString

func scannerWithParse[T any](parse parser[T]) scanner[T] {
	return func(vars Variables, name string, ptr *T) error {
		var str string
		if err := scanString(vars, name, &str); err != nil {
			return err
		}
		val, err := parse(str)
		if err != nil {
			return err
		}
		*ptr = val
		return nil
	}
}

func scanOrElse[T any](scanner scanner[T], defaultVal T) scanner[T] {
	return func(vars Variables, name string, ptr *T) error {
		if err := scanner(vars, name, ptr); err != nil {
			if isMissingEnvVarError(err) {
				*ptr = defaultVal
				return nil
			}
			return err
		}
		return nil
	}
}
