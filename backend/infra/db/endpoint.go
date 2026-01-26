package db

import (
	"bytes"
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
	DataSourceName() string
}

type FileEndpoint struct {
	Path   string
	Params *ParameterSet
}

var _ Endpoint = (*FileEndpoint)(nil)

func (e *FileEndpoint) DataSourceName() string {
	buf := new(bytes.Buffer)
	buf.WriteString("file:")
	buf.WriteString(e.Path)
	if qs := e.Params.query(); len(qs) > 0 {
		buf.WriteByte('?')
		buf.WriteString(qs.Encode())
	}
	return buf.String()
}
