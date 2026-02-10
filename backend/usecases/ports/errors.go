package ports

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"slices"

	"github.com/aereal/iter/seq"
	"github.com/aereal/mt"
	"github.com/aereal/nikki/backend/domain"
)

type ConvertMTEntryError struct {
	ArticleID         domain.ArticleID
	ArticleRevisionID domain.ArticleRevisionID
	Errs              []error
}

var _ error = (*ConvertMTEntryError)(nil)

func (e *ConvertMTEntryError) Error() string {
	buf := new(bytes.Buffer)
	buf.WriteString("failed to convert MT entry: ")
	var seen bool
	for _, err := range e.Errs {
		if seen {
			buf.WriteString("; ")
		}
		buf.WriteString(err.Error())
		seen = true
	}
	return buf.String()
}

func (e *ConvertMTEntryError) Unwrap() []error {
	if e == nil {
		return nil
	}
	return slices.Clone(e.Errs)
}

func (e *ConvertMTEntryError) Is(other error) bool {
	rhs := new(ConvertMTEntryError)
	if !errors.As(other, &rhs) {
		return false
	}
	if len(e.Errs) != len(rhs.Errs) {
		return false
	}
	if e.ArticleID != rhs.ArticleID || e.ArticleRevisionID != rhs.ArticleRevisionID {
		return false
	}
	leftErrs := slices.SortedFunc(
		slices.Values(e.Errs),
		cmpErr,
	)
	rightErrs := slices.SortedFunc(
		slices.Values(rhs.Errs),
		cmpErr,
	)
	for a, b := range seq.Zip(slices.Values(leftErrs), slices.Values(rightErrs)) {
		if !errors.Is(a, b) {
			return false
		}
	}
	return true
}

type UnsupportedConvertBreaksError struct {
	Value mt.ConvertBreaks
}

var _ error = (*UnsupportedConvertBreaksError)(nil)

func (e *UnsupportedConvertBreaksError) Error() string {
	return fmt.Sprintf("unsupported convert breaks: %s", e.Value)
}

func (e *UnsupportedConvertBreaksError) Is(other error) bool {
	rhs := new(UnsupportedConvertBreaksError)
	if !errors.As(other, &rhs) {
		return false
	}
	return e.Value == rhs.Value
}

var ErrEmptyBasename EmptyBasenameError

type EmptyBasenameError struct{}

func (EmptyBasenameError) Error() string { return "empty Basename" }

var ErrEmptyDate EmptyDateError

type EmptyDateError struct{}

func (EmptyDateError) Error() string { return "empty Date" }

func cmpErr(a, b error) int { return cmp.Compare(a.Error(), b.Error()) }
