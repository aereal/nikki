package types

import "time"

type DateTimeParser string

func (p DateTimeParser) Parse(str string) (time.Time, error) {
	parsed, err := time.Parse(string(p), str)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}

func (p DateTimeParser) Format(t time.Time) string { return t.Format(string(p)) }
