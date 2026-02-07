package ports

import "io"

type MTExport interface {
	OpenExport() (io.Reader, error)
}
