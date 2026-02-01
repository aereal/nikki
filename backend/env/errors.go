package env

type MissingEnvironmentVariableError struct {
	Name string
}

var _ error = (*MissingEnvironmentVariableError)(nil)

func (e *MissingEnvironmentVariableError) Error() string {
	return "missing environment variable: " + e.Name
}

func (e *MissingEnvironmentVariableError) Is(other error) bool {
	return isMissingEnvVarError(other)
}

func asMissingEnvVarError(err error) (*MissingEnvironmentVariableError, bool) {
	missingErr, ok := err.(*MissingEnvironmentVariableError)
	return missingErr, ok
}

func isMissingEnvVarError(err error) bool {
	_, ok := asMissingEnvVarError(err)
	return ok
}
