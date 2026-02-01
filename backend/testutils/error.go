package testutils

type LiteralError string

func (l LiteralError) Error() string { return string(l) }

func (l LiteralError) Is(other error) bool {
	if other == nil {
		return false
	}
	return string(l) == other.Error()
}
