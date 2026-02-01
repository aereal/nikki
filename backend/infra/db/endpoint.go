package db

import (
	"net/url"
)

type CacheMode int

const (
	CacheModeUnset CacheMode = iota
	CacheModeShared
	CacheModePrivate
)

type ParameterSet struct {
	Cache CacheMode
}

func (ps *ParameterSet) query() url.Values {
	q := url.Values{}
	if ps == nil {
		return q
	}
	switch ps.Cache {
	case CacheModeShared:
		q.Set("cache", "shared")
	case CacheModePrivate:
		q.Set("cache", "private")
	default: // noop
	}
	return q
}

type Endpoint interface {
	DataSourceName() (string, error)
}

type FileEndpoint struct {
	Path   string
	Params *ParameterSet
}

var _ Endpoint = (*FileEndpoint)(nil)

func (e *FileEndpoint) DataSourceName() (string, error) {
	if e.Path == "" {
		return "", ErrEmptyFile
	}

	parsed, err := url.Parse(e.Path)
	if err != nil {
		return "", err
	}
	parsed.Scheme = "file"
	if qs := e.Params.query(); len(qs) > 0 {
		parsed.RawQuery = qs.Encode()
	}
	return parsed.String(), nil
}
